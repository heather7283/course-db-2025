package main

import (
	"github.com/AllenDang/cimgui-go/imgui"
)

var uiState struct {
	countriesOpened bool
	countriesList []Country
	countryCodeInput string
	countryNameInput string

	hasError bool
	error string
}

func showCountries() {
	if !uiState.countriesOpened {
		uiState.countriesOpened = true
		uiState.countriesList, _ = getCountries()
	}

	avail := imgui.ContentRegionAvail()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##countryCodeInput", "Код страны", &uiState.countryCodeInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##countryNameInput", "Название страны", &uiState.countryNameInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 1)
	if imgui.Button("Добавить") {
		if err := addCountry(uiState.countryCodeInput, uiState.countryNameInput); err != nil {
			uiState.hasError = true
			uiState.error = err.Error()
		} else {
			uiState.countriesList, _ = getCountries()
		}
	}

	tableFlags := imgui.TableFlagsBordersOuter|imgui.TableFlagsBordersInner|imgui.TableFlagsRowBg|imgui.TableFlagsScrollY|imgui.TableFlagsScrollX;
	if imgui.BeginTableV("##countriesTable", 3, tableFlags, imgui.Vec2{}, 0) {
		defer imgui.EndTable()
		for i, c := range uiState.countriesList {
			imgui.PushIDInt(int32(i))

			imgui.TableNextRow()
			imgui.TableNextColumn()
			if (imgui.Button("x")) {
				deleteCountry(c.Code)
				uiState.countriesList, _ = getCountries()
			}
			imgui.TableNextColumn()
			imgui.Text(c.Code)
			imgui.TableNextColumn()
			imgui.Text(c.Name)

			imgui.PopID()
		}
	}
}

func runUI() {
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSize(imgui.CurrentIO().DisplaySize())
	imgui.BeginV("##mainWin", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoDecoration)


	if imgui.BeginTabBar("##tabBar") {
		if imgui.BeginTabItem("Страны") {
			showCountries()
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Виды спорта") {
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Атлеты") {
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Команды") {
			imgui.EndTabItem()
		}
		if imgui.BeginTabItem("Площадки") {
			imgui.EndTabItem()
		}

		imgui.EndTabBar()
	}

	// hack
	if uiState.hasError {
		imgui.OpenPopupStr("Ошибка")
		uiState.hasError = false
	}
	if imgui.BeginPopupModalV("Ошибка", nil, imgui.WindowFlagsNoResize) {
		imgui.TextUnformatted(uiState.error)
		if (imgui.Button("OK")) {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	imgui.End()
}

func initUI() {
	uiState.countriesOpened = false
}

