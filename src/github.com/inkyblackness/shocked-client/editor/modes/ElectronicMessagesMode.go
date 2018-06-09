package modes

import (
	"fmt"
	"os"
	"path"

	"github.com/inkyblackness/res/audio"
	"github.com/inkyblackness/res/audio/wav"
	"github.com/inkyblackness/shocked-client/editor/cmd"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

type textVariant uint32

const (
	textVariantVerbose = 0
	textVariantTerse   = 1
)

// duplicate to ElectronicMessages.go - so far no need to transport this.
var messageRanges = map[dataModel.ElectronicMessageType]int64{
	dataModel.ElectronicMessageTypeMail:     0x09B8 - 0x0989,
	dataModel.ElectronicMessageTypeLog:      0x0A98 - 0x09B8,
	dataModel.ElectronicMessageTypeFragment: 0x0AA8 - 0x0A98}

// ElectronicMessagesMode is a mode for messages.
type ElectronicMessagesMode struct {
	context        Context
	messageAdapter *model.ElectronicMessageAdapter

	area *ui.Area

	propertiesArea *ui.Area

	messageTypeLabel    *controls.Label
	messageTypeBox      *controls.ComboBox
	messageTypeItems    enumItems
	selectedMessageType dataModel.ElectronicMessageType
	messageTypeByIndex  map[uint32]dataModel.ElectronicMessageType
	messageIDLabel      *controls.Label
	messageIDSlider     *controls.Slider
	selectedMessageID   int

	removeLabel  *controls.Label
	removeButton *controls.TextButton

	propertiesHeader *controls.Label

	languageLabel    *controls.Label
	languageBox      *controls.ComboBox
	languageItems    enumItems
	selectedLanguage dataModel.ResourceLanguage

	variantLabel    *controls.Label
	variantBox      *controls.ComboBox
	variantItems    enumItems
	selectedVariant textVariant

	titleLabel        *controls.Label
	titleValue        *controls.Label
	nextMessageLabel  *controls.Label
	nextMessageValue  *controls.Slider
	isInterruptLabel  *controls.Label
	isInterruptBox    *controls.ComboBox
	isInterruptItems  map[bool]controls.ComboBoxItem
	colorLabel        *controls.Label
	colorValue        *controls.Slider
	leftDisplayLabel  *controls.Label
	leftDisplayValue  *controls.Slider
	rightDisplayLabel *controls.Label
	rightDisplayValue *controls.Slider

	audioArea       *ui.Area
	audioLabel      *controls.Label
	audioInfo       *controls.Label
	audioDropTarget *ui.Area

	displayArea *ui.Area

	leftDisplay  *controls.ImageDisplay
	rightDisplay *controls.ImageDisplay
	textValue    *controls.Label
	subjectValue *controls.Label
	senderValue  *controls.Label
}

// NewElectronicMessagesMode returns a new instance.
func NewElectronicMessagesMode(context Context, parent *ui.Area) *ElectronicMessagesMode {
	mode := &ElectronicMessagesMode{
		context:            context,
		messageAdapter:     context.ModelAdapter().ElectronicMessageAdapter(),
		messageTypeByIndex: make(map[uint32]dataModel.ElectronicMessageType),
		selectedLanguage:   dataModel.ResourceLanguageStandard,
		selectedMessageID:  -1,
		isInterruptItems:   make(map[bool]controls.ComboBoxItem)}

	scaled := func(value float32) float32 {
		return value * context.ControlFactory().Scale()
	}

	indexByMessageType := make(map[dataModel.ElectronicMessageType]uint32)
	for index, messageType := range dataModel.ElectronicMessageTypes() {
		mode.messageTypeByIndex[uint32(index)] = messageType
		indexByMessageType[messageType] = uint32(index)
	}

	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(parent)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(false)
		mode.area = builder.Build()
	}
	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewOffsetAnchor(parent.Top(), 0))
		builder.SetRight(ui.NewRelativeAnchor(parent.Left(), parent.Right(), 0.5))
		builder.SetBottom(ui.NewRelativeAnchor(parent.Top(), parent.Bottom(), 0.66))
		builder.SetVisible(true)
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.3))
		})
		builder.OnEvent(events.MouseMoveEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonDownEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)
		mode.propertiesArea = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.propertiesArea, context.ControlFactory())

		{
			mode.messageTypeLabel, mode.messageTypeBox = panelBuilder.addComboProperty("Message Type", mode.onMessageTypeChanged)
			mode.messageTypeItems = []*enumItem{
				{indexByMessageType[dataModel.ElectronicMessageTypeMail], "Mail"},
				{indexByMessageType[dataModel.ElectronicMessageTypeLog], "Log"},
				{indexByMessageType[dataModel.ElectronicMessageTypeFragment], "Fragment"}}
			mode.messageTypeBox.SetItems(mode.messageTypeItems.forComboBox())
		}
		{
			mode.messageIDLabel, mode.messageIDSlider = panelBuilder.addSliderProperty("Selected Message ID",
				func(newValue int64) { mode.onMessageSelected(int(newValue)) })
		}
		{
			mode.removeLabel, mode.removeButton = panelBuilder.addTextButton("Remove Selected", "Remove", mode.removeMessage)
		}
		mode.propertiesHeader = panelBuilder.addTitle("Properties")
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			mode.languageItems = []*enumItem{
				{uint32(dataModel.ResourceLanguageStandard), "STD"},
				{uint32(dataModel.ResourceLanguageFrench), "FRN"},
				{uint32(dataModel.ResourceLanguageGerman), "GER"}}
			mode.languageBox.SetItems(mode.languageItems.forComboBox())
		}
		{
			mode.variantLabel, mode.variantBox = panelBuilder.addComboProperty("Text Variant", mode.onVariantChanged)
			mode.variantItems = []*enumItem{{textVariantVerbose, "Verbose"}, {textVariantTerse, "Terse"}}
			mode.variantBox.SetItems(mode.variantItems.forComboBox())
		}
		mode.titleLabel, mode.titleValue = panelBuilder.addInfo("Title")
		mode.titleValue.AllowTextChange(mode.onTitleChangeRequested)
		mode.nextMessageLabel, mode.nextMessageValue = panelBuilder.addSliderProperty("Next Message", mode.onNextMessageChanged)
		mode.nextMessageValue.SetRange(-1, 0xFF)
		{
			mode.isInterruptLabel, mode.isInterruptBox = panelBuilder.addComboProperty("Is Interrupt", mode.onIsInterruptChanged)
			items := []controls.ComboBoxItem{
				&enumItem{0, "false"},
				&enumItem{1, "true"}}
			mode.isInterruptItems[false] = items[0]
			mode.isInterruptItems[true] = items[1]
			mode.isInterruptBox.SetItems(items)
			mode.isInterruptBox.SetSelectedItem(items[0])
		}
		mode.colorLabel, mode.colorValue = panelBuilder.addSliderProperty("Color Index", mode.onColorIndexChanged)
		mode.colorValue.SetRange(-1, 0xFF)
		mode.leftDisplayLabel, mode.leftDisplayValue = panelBuilder.addSliderProperty("Left Display", mode.onLeftDisplayChanged)
		mode.leftDisplayValue.SetRange(-1, 0xFF)
		mode.rightDisplayLabel, mode.rightDisplayValue = panelBuilder.addSliderProperty("Right Display", mode.onRightDisplayChanged)
		mode.rightDisplayValue.SetRange(-1, 0xFF)

		var audioBuilder *controlPanelBuilder
		mode.audioArea, audioBuilder = panelBuilder.addSection(false)
		mode.audioLabel, mode.audioInfo = audioBuilder.addInfo("Audio")
		audioDropTargetBuilder := ui.NewAreaBuilder()
		audioDropTargetBuilder.SetParent(mode.audioArea)
		audioDropTargetBuilder.SetLeft(ui.NewOffsetAnchor(mode.audioArea.Left(), 0))
		audioDropTargetBuilder.SetTop(ui.NewOffsetAnchor(mode.audioArea.Top(), 0))
		audioDropTargetBuilder.SetRight(ui.NewOffsetAnchor(mode.audioArea.Right(), 0))
		audioDropTargetBuilder.SetBottom(ui.NewOffsetAnchor(mode.audioArea.Bottom(), 0))
		audioDropTargetBuilder.OnEvent(events.FileDropEventType, mode.onAudioFileDropped)
		mode.audioDropTarget = audioDropTargetBuilder.Build()
	}
	{
		builder := ui.NewAreaBuilder()
		builder.SetParent(mode.area)
		builder.SetLeft(ui.NewOffsetAnchor(parent.Left(), 0))
		builder.SetTop(ui.NewRelativeAnchor(parent.Top(), parent.Bottom(), 0.66))
		builder.SetRight(ui.NewOffsetAnchor(parent.Right(), 0))
		builder.SetBottom(ui.NewOffsetAnchor(parent.Bottom(), 0))
		builder.SetVisible(true)
		builder.OnRender(func(area *ui.Area) {
			context.ForGraphics().RectangleRenderer().Fill(
				area.Left().Value(), area.Top().Value(), area.Right().Value(), area.Bottom().Value(),
				graphics.RGBA(0.7, 0.0, 0.7, 0.3))
		})
		builder.OnEvent(events.MouseMoveEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonUpEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonDownEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseButtonClickedEventType, ui.SilentConsumer)
		builder.OnEvent(events.MouseScrollEventType, ui.SilentConsumer)
		mode.displayArea = builder.Build()
	}
	padding := scaled(5)
	{
		labelBuilder := mode.context.ControlFactory().ForLabel()

		labelBuilder.SetParent(mode.displayArea)
		labelBuilder.SetTop(ui.NewOffsetAnchor(mode.displayArea.Top(), padding))
		labelBuilder.SetBottom(ui.NewOffsetAnchor(mode.displayArea.Bottom(), -padding))
		labelBuilder.SetLeft(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.25))
		labelBuilder.SetRight(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.75))
		labelBuilder.AlignedHorizontallyBy(controls.LeftAligner)
		labelBuilder.AlignedVerticallyBy(controls.LeftAligner)
		labelBuilder.SetFitToWidth()
		mode.textValue = labelBuilder.Build()
		mode.textValue.AllowTextChange(mode.onMessageTextChangeRequested)
	}
	{
		builder := mode.context.ControlFactory().ForImageDisplay()

		builder.SetParent(mode.displayArea)
		builder.SetTop(ui.NewOffsetAnchor(mode.displayArea.Top(), padding))
		builder.SetBottom(ui.NewOffsetAnchor(mode.displayArea.Bottom(), -padding))
		builder.SetLeft(ui.NewOffsetAnchor(mode.displayArea.Left(), padding))
		builder.SetRight(ui.NewOffsetAnchor(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.25), -padding))
		builder.WithProvider(mode.leftDisplayImage)
		mode.leftDisplay = builder.Build()
	}
	{
		builder := mode.context.ControlFactory().ForImageDisplay()

		builder.SetParent(mode.displayArea)
		builder.SetTop(ui.NewOffsetAnchor(mode.displayArea.Top(), padding))
		builder.SetBottom(ui.NewOffsetAnchor(mode.displayArea.Bottom(), -padding))
		builder.SetLeft(ui.NewOffsetAnchor(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.75), padding))
		builder.SetRight(ui.NewOffsetAnchor(mode.displayArea.Right(), -padding))
		builder.WithProvider(mode.rightDisplayImage)
		mode.rightDisplay = builder.Build()
	}
	{
		labelBuilder := mode.context.ControlFactory().ForLabel()

		labelBuilder.SetParent(mode.displayArea)
		labelBuilder.SetTop(ui.NewOffsetAnchor(mode.displayArea.Top(), padding))
		labelBuilder.SetBottom(ui.NewRelativeAnchor(mode.displayArea.Top(), mode.displayArea.Bottom(), 0.5))
		labelBuilder.SetLeft(ui.NewOffsetAnchor(mode.displayArea.Left(), padding))
		labelBuilder.SetRight(ui.NewOffsetAnchor(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.25), -padding))
		labelBuilder.AlignedHorizontallyBy(controls.LeftAligner)
		labelBuilder.AlignedVerticallyBy(controls.LeftAligner)
		labelBuilder.SetFitToWidth()
		mode.senderValue = labelBuilder.Build()
		mode.senderValue.AllowTextChange(mode.onSenderChangeRequested)
	}
	{
		labelBuilder := mode.context.ControlFactory().ForLabel()

		labelBuilder.SetParent(mode.displayArea)
		labelBuilder.SetTop(ui.NewRelativeAnchor(mode.displayArea.Top(), mode.displayArea.Bottom(), 0.5))
		labelBuilder.SetBottom(ui.NewOffsetAnchor(mode.displayArea.Bottom(), -padding))
		labelBuilder.SetLeft(ui.NewOffsetAnchor(mode.displayArea.Left(), padding))
		labelBuilder.SetRight(ui.NewOffsetAnchor(ui.NewRelativeAnchor(mode.displayArea.Left(), mode.displayArea.Right(), 0.25), -padding))
		labelBuilder.AlignedHorizontallyBy(controls.LeftAligner)
		labelBuilder.AlignedVerticallyBy(controls.RightAligner)
		labelBuilder.SetFitToWidth()
		mode.subjectValue = labelBuilder.Build()
		mode.subjectValue.AllowTextChange(mode.onSubjectChangeRequested)
	}
	mode.messageAdapter.OnMessageDataChanged(mode.onMessageDataChanged)
	mode.messageAdapter.OnMessageAudioChanged(mode.onMessageAudioChanged)

	mode.setState(dataModel.ElectronicMessageTypeMail, 0, dataModel.ResourceLanguageStandard, textVariantVerbose)
	mode.context.ModelAdapter().OnProjectChanged(func() {
		mode.requestData()
	})

	return mode
}

// SetActive implements the Mode interface.
func (mode *ElectronicMessagesMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *ElectronicMessagesMode) leftDisplayImage() (texture *graphics.BitmapTexture) {
	return mode.displayImage(mode.messageAdapter.LeftDisplay())
}

func (mode *ElectronicMessagesMode) rightDisplayImage() (texture *graphics.BitmapTexture) {
	return mode.displayImage(mode.messageAdapter.RightDisplay())
}

func (mode *ElectronicMessagesMode) displayImage(index int) (texture *graphics.BitmapTexture) {
	if (index >= 0) && (index < 0x100) {
		resourceKey := dataModel.MakeLocalizedResourceKey(dataModel.ResourceTypeMfdDataImages, mode.selectedLanguage, uint16(index))
		texture = mode.context.ForGraphics().BitmapsStore().Texture(graphics.TextureKeyFromInt(resourceKey.ToInt()))
	}
	return
}

func (mode *ElectronicMessagesMode) onAudioFileDropped(area *ui.Area, event events.Event) (consumed bool) {
	dropEvent := event.(*events.FileDropEvent)

	if len(dropEvent.FilePaths()) == 1 {
		filePath := dropEvent.FilePaths()[0]
		fileInfo, err := os.Stat(filePath)

		if err == nil {
			if fileInfo.IsDir() {
				mode.exportAudio(filePath)
			} else {
				mode.importAudio(filePath)
			}
		} else {
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File is not found/recognized %s", filePath))
		}
		consumed = true
	}

	return
}

func (mode *ElectronicMessagesMode) exportAudio(filePath string) {
	languageIndex := mode.selectedLanguage.ToIndex()
	soundData := mode.messageAdapter.Audio(languageIndex)

	if soundData != nil {
		fileName := path.Join(filePath, fmt.Sprintf("%v_%02d_%v.wav", mode.selectedMessageType, mode.selectedMessageID, mode.selectedLanguage.ShortName()))
		file, err := os.Create(fileName)

		if err == nil {
			defer func() {
				_ = file.Close()
			}()
			wav.Save(file, soundData.SampleRate(), soundData.Samples(0, soundData.SampleCount()))
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Exported %s", fileName))
		} else {
			mode.context.ModelAdapter().SetMessage("Could not create file for export.")
		}
	}
}

func (mode *ElectronicMessagesMode) importAudio(filePath string) {
	file, fileErr := os.Open(filePath)

	if (fileErr == nil) && (file != nil) {
		defer func() {
			_ = file.Close()
		}()
		data, dataErr := wav.Load(file)

		if dataErr == nil {
			mode.requestAudioChange(data)
		} else {
			mode.context.ModelAdapter().SetMessage("File not supported. Only .wav files with 16bit or 8bit LPCM possible.")
		}
	} else {
		mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File could not be opened: %s", filePath))
	}
}

func (mode *ElectronicMessagesMode) requestAudioChange(newData audio.SoundData) {
	restoreState := mode.stateSnapshot()
	mode.context.Perform(&cmd.SetAudioCommand{
		Setter: func(data audio.SoundData) error {
			restoreState()
			mode.messageAdapter.RequestAudioChange(mode.selectedLanguage, data)
			return nil
		},
		NewValue: newData,
		OldValue: mode.messageAdapter.Audio(mode.selectedLanguage.ToIndex())})
}

func (mode *ElectronicMessagesMode) onMessageTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.setState(mode.messageTypeByIndex[item.value], 0, mode.selectedLanguage, mode.selectedVariant)
}

func (mode *ElectronicMessagesMode) onMessageSelected(id int) {
	mode.selectedMessageID = id
	mode.requestData()
}

func (mode *ElectronicMessagesMode) requestData() {
	mode.messageAdapter.RequestMessage(mode.selectedMessageType, mode.selectedMessageID)
}

func (mode *ElectronicMessagesMode) removeMessage() {
	restoreState := mode.stateSnapshot()
	command := &cmd.RemoveElectronicMessageCommand{
		RestoreState: restoreState,
		Store:        mode.messageAdapter}

	command.Properties.NextMessage = intAsPointer(mode.messageAdapter.NextMessage())
	command.Properties.IsInterrupt = boolAsPointer(mode.messageAdapter.IsInterrupt())
	command.Properties.ColorIndex = intAsPointer(mode.messageAdapter.ColorIndex())
	command.Properties.LeftDisplay = intAsPointer(mode.messageAdapter.LeftDisplay())
	command.Properties.RightDisplay = intAsPointer(mode.messageAdapter.RightDisplay())

	for langIndex := 0; langIndex < dataModel.LanguageCount; langIndex++ {
		command.Properties.Subject[langIndex] = stringAsPointer(mode.messageAdapter.Subject(langIndex))
		command.Properties.Sender[langIndex] = stringAsPointer(mode.messageAdapter.Sender(langIndex))
		command.Properties.Title[langIndex] = stringAsPointer(mode.messageAdapter.Title(langIndex))
		command.Properties.VerboseText[langIndex] = stringAsPointer(mode.messageAdapter.VerboseText(langIndex))
		command.Properties.TerseText[langIndex] = stringAsPointer(mode.messageAdapter.TerseText(langIndex))

		command.Audio[langIndex] = mode.messageAdapter.Audio(langIndex)
	}

	mode.context.Perform(command)
}

func (mode *ElectronicMessagesMode) onMessageDataChanged() {
	mode.updateMessageText()
	mode.updateMessageData()
}

func (mode *ElectronicMessagesMode) onMessageAudioChanged() {
	mode.updateMessageAudio()
}

func (mode *ElectronicMessagesMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedLanguage = dataModel.ResourceLanguage(item.value)
	mode.updateMessageText()
	mode.updateMessageAudio()
}

func (mode *ElectronicMessagesMode) onVariantChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedVariant = textVariant(item.value)
	mode.updateMessageText()
}

func (mode *ElectronicMessagesMode) updateMessageText() {
	languageIndex := mode.selectedLanguage.ToIndex()
	var text string
	if mode.selectedVariant == textVariantTerse {
		text = mode.messageAdapter.TerseText(languageIndex)
	} else {
		text = mode.messageAdapter.VerboseText(languageIndex)
	}
	mode.textValue.SetText(text)

	mode.subjectValue.SetText(mode.messageAdapter.Subject(languageIndex))
	mode.titleValue.SetText(mode.messageAdapter.Title(languageIndex))
	mode.senderValue.SetText(mode.messageAdapter.Sender(languageIndex))
}

func (mode *ElectronicMessagesMode) updateMessageData() {
	mode.nextMessageValue.SetValue(int64(mode.messageAdapter.NextMessage()))
	mode.isInterruptBox.SetSelectedItem(mode.isInterruptItems[mode.messageAdapter.IsInterrupt()])
	mode.colorValue.SetValue(int64(mode.messageAdapter.ColorIndex()))
	mode.leftDisplayValue.SetValue(int64(mode.messageAdapter.LeftDisplay()))
	mode.rightDisplayValue.SetValue(int64(mode.messageAdapter.RightDisplay()))
}

func (mode *ElectronicMessagesMode) updateMessageAudio() {
	languageIndex := mode.selectedLanguage.ToIndex()
	data := mode.messageAdapter.Audio(languageIndex)
	var info string

	if data != nil {
		info = fmt.Sprintf("%.02f sec", float32(data.SampleCount())/data.SampleRate())
	} else {
		info = "(no audio)"
	}

	mode.audioInfo.SetText(info)
}

func (mode *ElectronicMessagesMode) onNextMessageChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.ElectronicMessage, value *int) {
		properties.NextMessage = value
	}, int(newValue), mode.messageAdapter.NextMessage())
}

func (mode *ElectronicMessagesMode) onIsInterruptChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	newValue := item.value != 0
	mode.requestBooleanPropertyChange(func(properties *dataModel.ElectronicMessage, value *bool) {
		properties.IsInterrupt = value
	}, newValue, mode.messageAdapter.IsInterrupt())
}

func (mode *ElectronicMessagesMode) onColorIndexChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.ElectronicMessage, value *int) {
		properties.ColorIndex = value
	}, int(newValue), mode.messageAdapter.ColorIndex())
}

func (mode *ElectronicMessagesMode) onLeftDisplayChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.ElectronicMessage, value *int) {
		properties.LeftDisplay = value
	}, int(newValue), mode.messageAdapter.LeftDisplay())
}

func (mode *ElectronicMessagesMode) onRightDisplayChanged(newValue int64) {
	mode.requestIntPropertyChange(func(properties *dataModel.ElectronicMessage, value *int) {
		properties.RightDisplay = value
	}, int(newValue), mode.messageAdapter.RightDisplay())
}

func (mode *ElectronicMessagesMode) onMessageTextChangeRequested(newText string) {
	var oldText string
	languageIndex := mode.selectedLanguage.ToIndex()
	if mode.selectedVariant == textVariantTerse {
		oldText = mode.messageAdapter.TerseText(languageIndex)
	} else {
		oldText = mode.messageAdapter.VerboseText(languageIndex)
	}
	mode.requestStringPropertyChange(func(properties *dataModel.ElectronicMessage, value *string) {
		if mode.selectedVariant == textVariantTerse {
			properties.TerseText[languageIndex] = value
		} else {
			properties.VerboseText[languageIndex] = value
		}
	}, newText, oldText)
}

func (mode *ElectronicMessagesMode) onSubjectChangeRequested(newText string) {
	languageIndex := mode.selectedLanguage.ToIndex()
	mode.requestStringPropertyChange(func(properties *dataModel.ElectronicMessage, value *string) {
		properties.Subject[languageIndex] = value
	}, newText, mode.messageAdapter.Subject(languageIndex))
}

func (mode *ElectronicMessagesMode) onSenderChangeRequested(newText string) {
	languageIndex := mode.selectedLanguage.ToIndex()
	mode.requestStringPropertyChange(func(properties *dataModel.ElectronicMessage, value *string) {
		properties.Sender[languageIndex] = value
	}, newText, mode.messageAdapter.Sender(languageIndex))
}

func (mode *ElectronicMessagesMode) onTitleChangeRequested(newText string) {
	languageIndex := mode.selectedLanguage.ToIndex()
	mode.requestStringPropertyChange(func(properties *dataModel.ElectronicMessage, value *string) {
		properties.Title[languageIndex] = value
	}, newText, mode.messageAdapter.Title(languageIndex))
}

func (mode *ElectronicMessagesMode) requestStringPropertyChange(modifier func(*dataModel.ElectronicMessage, *string),
	newValue, oldValue string) {
	restoreState := mode.stateSnapshot()

	mode.context.Perform(&cmd.SetStringPropertyCommand{
		Setter: func(value string) error {
			return mode.requestPropertyChange(restoreState, func(properties *dataModel.ElectronicMessage) { modifier(properties, &value) })
		},
		NewValue: newValue,
		OldValue: oldValue})
}

func (mode *ElectronicMessagesMode) requestBooleanPropertyChange(modifier func(*dataModel.ElectronicMessage, *bool),
	newValue, oldValue bool) {
	restoreState := mode.stateSnapshot()

	mode.context.Perform(&cmd.SetBooleanPropertyCommand{
		Setter: func(value bool) error {
			return mode.requestPropertyChange(restoreState, func(properties *dataModel.ElectronicMessage) { modifier(properties, &value) })
		},
		NewValue: newValue,
		OldValue: oldValue})
}

func (mode *ElectronicMessagesMode) requestIntPropertyChange(modifier func(*dataModel.ElectronicMessage, *int),
	newValue, oldValue int) {
	restoreState := mode.stateSnapshot()

	mode.context.Perform(&cmd.SetIntPropertyCommand{
		Setter: func(value int) error {
			return mode.requestPropertyChange(restoreState, func(properties *dataModel.ElectronicMessage) { modifier(properties, &value) })
		},
		NewValue: newValue,
		OldValue: oldValue})
}

func (mode *ElectronicMessagesMode) requestPropertyChange(restoreState func(), modifier func(*dataModel.ElectronicMessage)) error {
	restoreState()
	var properties dataModel.ElectronicMessage
	modifier(&properties)
	mode.messageAdapter.RequestMessageChange(properties)
	return nil
}

func (mode *ElectronicMessagesMode) stateSnapshot() func() {
	currentType := mode.selectedMessageType
	currentID := mode.selectedMessageID
	currentLanguage := mode.selectedLanguage
	currentVariant := mode.selectedVariant
	return func() {
		mode.setState(currentType, currentID, currentLanguage, currentVariant)
	}
}

func (mode *ElectronicMessagesMode) setState(messageType dataModel.ElectronicMessageType, id int,
	language dataModel.ResourceLanguage, variant textVariant) {
	{
		mode.selectedMessageType = messageType
		for _, item := range mode.messageTypeItems {
			if mode.messageTypeByIndex[item.value] == messageType {
				mode.messageTypeBox.SetSelectedItem(item)
			}
		}
		mode.audioArea.SetVisible(mode.selectedMessageType != dataModel.ElectronicMessageTypeFragment)
		mode.messageIDSlider.SetRange(0, messageRanges[mode.selectedMessageType]-1)
	}
	{
		mode.selectedLanguage = language
		for _, item := range mode.languageItems {
			if item.value == uint32(language) {
				mode.languageBox.SetSelectedItem(item)
			}
		}
	}
	{
		mode.selectedVariant = variant
		for _, item := range mode.variantItems {
			if item.value == uint32(variant) {
				mode.variantBox.SetSelectedItem(item)
			}
		}
	}
	mode.selectedMessageID = id
	mode.messageIDSlider.SetValue(int64(mode.selectedMessageID))
	mode.requestData()
}
