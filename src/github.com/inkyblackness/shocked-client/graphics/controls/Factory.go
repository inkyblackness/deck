package controls

// Factory is an interface for creating controls with a common look-and-feel.
type Factory interface {
	Scale() float32
	ForLabel() *LabelBuilder
	ForTextButton() *TextButtonBuilder
	ForComboBox() *ComboBoxBuilder
	ForTextureSelector() *TextureSelectorBuilder
	ForSlider() *SliderBuilder
	ForImageDisplay() *ImageDisplayBuilder
}
