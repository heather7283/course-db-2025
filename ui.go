package main

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
)

type Tab int

const (
	TabCountries Tab = iota
	TabSports Tab = iota
	TabAthletes Tab = iota
)

var tabs []Tab = []Tab{TabCountries, TabSports, TabAthletes}

var uiState struct {
	oldTab Tab

	countriesList []Country
	countryCodeInput string
	countryNameInput string

	sportsList []Sport
	sportCodeInput string
	sportNameInput string
	sportIsTeamInput bool

	athletesList []Athlete
	athleteNameInput string
	athleteIsMaleInput bool

	hasError bool
	error string
}

var tableFlags imgui.TableFlags =
	imgui.TableFlagsBordersOuter|imgui.TableFlagsBordersInner|imgui.TableFlagsRowBg|
	imgui.TableFlagsScrollY|imgui.TableFlagsScrollX|imgui.TableFlagsResizable|
	imgui.TableFlagsHighlightHoveredColumn;

func (tab Tab) name() string {
	switch tab {
	case TabCountries: return "Страны"
	case TabSports: return "Виды спорта"
	case TabAthletes: return "Спорстмены"
	default: return "INVALID TAB"
	}
}

func (tab Tab) show() {
	switched := false
	if uiState.oldTab != tab {
		switched = true
		uiState.oldTab = tab
	}

	switch tab {
	case TabCountries: showCountries(switched)
	case TabSports: showSports(switched)
	case TabAthletes: showAthletes(switched)
	default: showError(fmt.Errorf("INVALID TAB"))
	}
}

func showError(err error) {
	uiState.hasError = true
	uiState.error = err.Error()
}

func showAthletes(switched bool) {
	if switched {
		uiState.athletesList, _ = getAthletes()
	}

	//avail := imgui.ContentRegionAvail()
	//imgui.SetNextItemWidth(avail.X / 4)
	//imgui.InputTextWithHint("##athleteNameInput", "Имя", &uiState.athleteNameInput, 0, nil)
	//imgui.SameLine()
	//if imgui.RadioButtonBool("М", uiState.athleteIsMaleInput) {
	//	uiState.athleteIsMaleInput = true
	//} else if imgui.RadioButtonBool("Ж", !uiState.athleteIsMaleInput) {
	//	uiState.athleteIsMaleInput = false
	//}
	//imgui.SetNextItemWidth(avail.X / 4)
	//imgui.InputTextWithHint("##athleteBirthdayInput", "День рождения", &uiState.sportNameInput, 0, nil)
	//imgui.SameLine()
	//imgui.SetNextItemWidth(avail.X / 1)
	//if imgui.Button("Добавить") {
	//	err := addAthlete(uiState.athleteNameInput, uiState.athleteIsMaleInput, 0, )
	//	if err != nil {
	//		showError(err)
	//	} else {
	//		uiState.sportsList, _ = getSports()
	//	}
	//}

	if imgui.BeginTableV("##athletesTable", 5, tableFlags, imgui.Vec2{}, 0) {
		imgui.TableSetupColumn("")
		imgui.TableSetupColumn("Имя")
		imgui.TableSetupColumn("Пол")
		imgui.TableSetupColumn("День рождения")
		imgui.TableSetupColumn("Страна")
		imgui.TableHeadersRow()
		for i, a := range uiState.athletesList {
			var gender string
			if a.Gender == "M" {
				gender = "М"
			} else {
				gender = "Ж"
			}
			imgui.PushIDInt(int32(i))

			imgui.TableNextRow()
			imgui.TableNextColumn()
			if (imgui.Button("x")) {
				deleteAthlete(a.ID)
				uiState.athletesList, _ = getAthletes()
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.Name)
			imgui.TableNextColumn()
			imgui.TextUnformatted(gender)
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.Birthday.String())
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.CountryName)

			imgui.PopID()
		}
		imgui.EndTable()
	}
}

func showSports(switched bool) {
	if switched {
		uiState.sportsList, _ = getSports()
	}

	avail := imgui.ContentRegionAvail()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##sportCodeInput", "Код вида спорта", &uiState.sportCodeInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##sportNameInput", "Название вида спорта", &uiState.sportNameInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.Checkbox("Командный", &uiState.sportIsTeamInput)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 1)
	if imgui.Button("Добавить") {
		if err := addSport(uiState.countryCodeInput, uiState.countryNameInput, uiState.sportIsTeamInput); err != nil {
			showError(err)
		} else {
			uiState.sportsList, _ = getSports()
		}
	}

	if imgui.BeginTableV("##sportsTable", 4, tableFlags, imgui.Vec2{}, 0) {
		imgui.TableSetupColumn("")
		imgui.TableSetupColumn("Код")
		imgui.TableSetupColumn("Название")
		imgui.TableSetupColumn("Командный/одиночный")
		imgui.TableHeadersRow()
		for i, s := range uiState.sportsList {
			imgui.PushIDInt(int32(i))

			imgui.TableNextRow()
			imgui.TableNextColumn()
			if (imgui.Button("x")) {
				deleteSport(s.Code)
				uiState.sportsList, _ = getSports()
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(s.Code)
			imgui.TableNextColumn()
			imgui.TextUnformatted(s.Name)
			imgui.TableNextColumn()
			switch s.IsTeam{
			case true: imgui.TextUnformatted("Командный")
			case false: imgui.TextUnformatted("Одиночный")
			}

			imgui.PopID()
		}
		imgui.EndTable()
	}
}

func showCountries(switched bool) {
	if switched {
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
			showError(err)
		} else {
			uiState.countriesList, _ = getCountries()
		}
	}

	if imgui.BeginTableV("##countriesTable", 3, tableFlags, imgui.Vec2{}, 0) {
		imgui.TableSetupColumn("")
		imgui.TableSetupColumn("Код")
		imgui.TableSetupColumn("Название")
		imgui.TableHeadersRow()
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
		imgui.EndTable()
	}
}

func runUI() {
	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSize(imgui.CurrentIO().DisplaySize())
	imgui.BeginV("##mainWin", nil, imgui.WindowFlagsNoMove|imgui.WindowFlagsNoDecoration)


	if imgui.BeginTabBar("##tabBar") {
		for _, tab := range tabs {
			if imgui.BeginTabItem(tab.name()) {
				tab.show()
				imgui.EndTabItem()
			}
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
	uiState.oldTab = 100500
}

