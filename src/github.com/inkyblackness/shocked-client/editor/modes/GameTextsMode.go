package modes

import (
	"fmt"
	"os"
	"path"

	"github.com/inkyblackness/res/audio/wav"
	"github.com/inkyblackness/shocked-client/editor/model"
	"github.com/inkyblackness/shocked-client/graphics"
	"github.com/inkyblackness/shocked-client/graphics/controls"
	"github.com/inkyblackness/shocked-client/ui"
	"github.com/inkyblackness/shocked-client/ui/events"

	dataModel "github.com/inkyblackness/shocked-model"
)

// GameTextsMode is a mode for arbitrary game texts.
type GameTextsMode struct {
	context      Context
	textAdapter  *model.TextAdapter
	soundAdapter *model.SoundAdapter

	area           *ui.Area
	propertiesArea *ui.Area

	typeLabel            *controls.Label
	typeBox              *controls.ComboBox
	selectedResourceType dataModel.ResourceType

	languageLabel *controls.Label
	languageBox   *controls.ComboBox
	language      dataModel.ResourceLanguage

	selectedTextIDLabel  *controls.Label
	selectedTextIDSlider *controls.Slider
	selectedTextID       int

	textDrop  *ui.Area
	textValue *controls.Label

	audioArea       *ui.Area
	audioLabel      *controls.Label
	audioInfo       *controls.Label
	audioDropTarget *ui.Area
}

// NewGameTextsMode returns a new instance.
func NewGameTextsMode(context Context, parent *ui.Area) *GameTextsMode {
	mode := &GameTextsMode{
		context:        context,
		textAdapter:    context.ModelAdapter().TextAdapter(),
		soundAdapter:   context.ModelAdapter().SoundAdapter(),
		selectedTextID: -1}

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
		mode.propertiesArea = builder.Build()
	}
	{
		panelBuilder := newControlPanelBuilder(mode.propertiesArea, context.ControlFactory())
		var initialTypeItem controls.ComboBoxItem
		var initialLanguageItem controls.ComboBoxItem

		{
			mode.typeLabel, mode.typeBox = panelBuilder.addComboProperty("Text Type", mode.onTextTypeChanged)
			items := []controls.ComboBoxItem{
				&enumItem{uint32(dataModel.ResourceTypeTrapMessages), "Trap Messages"},
				&enumItem{uint32(dataModel.ResourceTypeWords), "Words"},
				&enumItem{uint32(dataModel.ResourceTypeLogCategories), "Log Categories"},
				&enumItem{uint32(dataModel.ResourceTypeVariousMessages), "Various Messages"},
				&enumItem{uint32(dataModel.ResourceTypeScreenMessages), "Screen Messages"},
				&enumItem{uint32(dataModel.ResourceTypeInfoNodeMessages), "Info Node Messages (8/5/6)"},
				&enumItem{uint32(dataModel.ResourceTypeAccessCardNames), "Access Card Names"},
				&enumItem{uint32(dataModel.ResourceTypeDataletMessages), "Datalet Messages (8/5/8)"},
				&enumItem{uint32(dataModel.ResourceTypePaperTexts), "Paper Texts"},
				&enumItem{uint32(dataModel.ResourceTypePanelNames), "Panel Names"}}

			mode.typeBox.SetItems(items)
			initialTypeItem = items[0]
		}
		{
			mode.languageLabel, mode.languageBox = panelBuilder.addComboProperty("Language", mode.onLanguageChanged)
			items := []controls.ComboBoxItem{
				&enumItem{uint32(dataModel.ResourceLanguageStandard), "STD"},
				&enumItem{uint32(dataModel.ResourceLanguageFrench), "FRN"},
				&enumItem{uint32(dataModel.ResourceLanguageGerman), "GER"}}
			mode.languageBox.SetItems(items)
			initialLanguageItem = items[0]
		}
		{
			mode.selectedTextIDLabel, mode.selectedTextIDSlider = panelBuilder.addSliderProperty("Selected Text ID",
				func(newValue int64) {
					mode.onTextSelected(int(newValue))
				})
		}
		{
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

			mode.soundAdapter.OnAudioChanged(mode.onAudioChanged)
		}
		mode.languageBox.SetSelectedItem(initialLanguageItem)
		mode.onLanguageChanged(initialLanguageItem)
		mode.typeBox.SetSelectedItem(initialTypeItem)
		mode.onTextTypeChanged(initialTypeItem)
	}
	{
		padding := float32(5.0)

		{
			dropBuilder := ui.NewAreaBuilder()
			displayBuilder := mode.context.ControlFactory().ForLabel()
			left := ui.NewOffsetAnchor(mode.propertiesArea.Right(), padding)
			right := ui.NewOffsetAnchor(mode.area.Right(), -padding)
			top := ui.NewOffsetAnchor(mode.area.Top(), padding)
			bottom := ui.NewOffsetAnchor(mode.area.Bottom(), -padding)

			dropBuilder.SetParent(mode.area)
			displayBuilder.SetParent(mode.area)
			dropBuilder.SetLeft(left)
			displayBuilder.SetLeft(left)
			dropBuilder.SetRight(right)
			displayBuilder.SetRight(right)
			dropBuilder.SetTop(top)
			displayBuilder.SetTop(top)
			dropBuilder.SetBottom(bottom)
			displayBuilder.SetBottom(bottom)
			displayBuilder.AlignedHorizontallyBy(controls.LeftAligner)
			displayBuilder.AlignedVerticallyBy(controls.LeftAligner)
			displayBuilder.SetFitToWidth()
			mode.textDrop = dropBuilder.Build()
			mode.textValue = displayBuilder.Build()
			mode.textValue.AllowTextChange(mode.onTextModified)
		}
	}
	mode.context.ModelAdapter().OnProjectChanged(func() {
		mode.onTextSelected(mode.selectedTextID)
	})
	mode.textAdapter.OnTextChanged(mode.onTextChanged)

	return mode
}

// SetActive implements the Mode interface.
func (mode *GameTextsMode) SetActive(active bool) {
	mode.area.SetVisible(active)
}

func (mode *GameTextsMode) onTextTypeChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.selectedResourceType = dataModel.ResourceType(item.value)

	mode.audioArea.SetVisible(mode.selectedResourceType == dataModel.ResourceTypeTrapMessages)
	mode.onTextSelected(0)
	mode.selectedTextIDSlider.SetRange(0, int64(dataModel.MaxEntriesFor(mode.selectedResourceType))-1)
	mode.selectedTextIDSlider.SetValue(0)
	mode.requestData()
}

func (mode *GameTextsMode) onLanguageChanged(boxItem controls.ComboBoxItem) {
	item := boxItem.(*enumItem)
	mode.language = dataModel.ResourceLanguage(item.value)
	mode.requestData()
}

func (mode *GameTextsMode) onTextSelected(id int) {
	mode.selectedTextID = id
	mode.requestData()
}

func (mode *GameTextsMode) requestData() {
	key := dataModel.MakeLocalizedResourceKey(mode.selectedResourceType, mode.language, uint16(mode.selectedTextID))
	mode.textAdapter.RequestText(key)

	if mode.selectedResourceType == dataModel.ResourceTypeTrapMessages {
		mode.requestAudio(dataModel.ResourceTypeTrapAudio)
	}
}

func (mode *GameTextsMode) onTextChanged() {
	mode.textValue.SetText(mode.textAdapter.Text())
}

func (mode *GameTextsMode) onTextModified(newText string) {
	mode.textAdapter.RequestTextChange(newText)
}

func (mode *GameTextsMode) requestAudio(resourceType dataModel.ResourceType) {
	key := dataModel.MakeLocalizedResourceKey(resourceType, mode.language, uint16(mode.selectedTextID))
	mode.context.ModelAdapter().SoundAdapter().RequestAudio(key)
}

func (mode *GameTextsMode) onAudioChanged() {
	data := mode.soundAdapter.Audio()
	info := ""

	if data != nil {
		info = fmt.Sprintf("%.02f sec", float32(data.SampleCount())/data.SampleRate())
	} else {
		info = "(no audio)"
	}

	mode.audioInfo.SetText(info)
}

func (mode *GameTextsMode) onAudioFileDropped(area *ui.Area, event events.Event) (consumed bool) {
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

func (mode *GameTextsMode) exportAudio(filePath string) {
	soundData := mode.soundAdapter.Audio()

	if soundData != nil {
		fileName := path.Join(filePath, fmt.Sprintf("traps_%02d_%v.wav", mode.selectedTextID, mode.language.ShortName()))
		file, err := os.Create(fileName)

		if err == nil {
			defer file.Close()
			wav.Save(file, soundData.SampleRate(), soundData.Samples(0, soundData.SampleCount()))
			mode.context.ModelAdapter().SetMessage(fmt.Sprintf("Exported %s", fileName))
		} else {
			mode.context.ModelAdapter().SetMessage("Could not create file for export.")
		}
	}
}

func (mode *GameTextsMode) importAudio(filePath string) {
	file, fileErr := os.Open(filePath)

	if (fileErr == nil) && (file != nil) {
		defer file.Close()
		data, dataErr := wav.Load(file)

		if dataErr == nil {
			mode.soundAdapter.RequestAudioChange(data)
		} else {
			mode.context.ModelAdapter().SetMessage("File not supported. Only .wav files with 16bit or 8bit LPCM possible.")
		}
	} else {
		mode.context.ModelAdapter().SetMessage(fmt.Sprintf("File could not be opened: %s", filePath))
	}
}
