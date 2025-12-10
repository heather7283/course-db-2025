package main

import (
	"fmt"
	"sort"
	"strings"

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

	athletesDirty bool
	athletesSortSwitch int32
	athletesSortFunc func(a *Athlete, b *Athlete) bool
	athletesList []Athlete
	athletesListProcessed []*Athlete
	athleteNameInput string
	athleteIsMaleInput bool
	athleteNameFilter string
	athleteGenderFilter int32

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

func showTable[T any](id string, headers []string, items []T, callback func(item T)) {
	cols := len(headers)
	if imgui.BeginTableV(id, int32(cols), tableFlags, imgui.Vec2{}, 0) {
		imgui.TableSetupScrollFreeze(0, 1)

		for _, header := range headers {
			imgui.TableSetupColumn(header)
		}
		imgui.TableHeadersRow()

		for i := range items {
			imgui.PushIDInt(int32(i))
			callback(items[i])
			imgui.PopID()
		}
		imgui.EndTable()
	}
}

func processAthletes() {
	uiState.athletesListProcessed = make([]*Athlete, 0, len(uiState.athletesList))

	for i := range uiState.athletesList {
		a := &uiState.athletesList[i]
		if (len(uiState.athleteNameFilter) > 0 && !strings.Contains(a.Name, uiState.athleteNameFilter)) {
			continue
		} else if (uiState.athleteGenderFilter > 0) {
			if (uiState.athleteGenderFilter == 1 && a.Gender != "M" ||
				uiState.athleteGenderFilter == 2 && a.Gender != "F") {
				continue
			}
		}
		uiState.athletesListProcessed = append(uiState.athletesListProcessed, a)
	}

	if uiState.athletesSortFunc != nil {
		sort.Slice(uiState.athletesListProcessed, func(i, j int) bool {
			return uiState.athletesSortFunc(uiState.athletesListProcessed[i], uiState.athletesListProcessed[j])
		})
	}
}

func showAthletes(switched bool) {
	if switched {
		uiState.athletesList, _ = getAthletes()
		uiState.athletesDirty = true
	}

	// TODO
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

	avail := imgui.ContentRegionAvail()
	imgui.TextUnformatted("Фильтр")

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##athletesFilterName", "Имя", &uiState.athleteNameFilter, 0, nil) {
		uiState.athletesDirty = true
	}

	if imgui.SameLine(); imgui.RadioButtonIntPtr("Все##athletesFilter", &uiState.athleteGenderFilter, 0) {
		uiState.athletesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("М##athletesFilter", &uiState.athleteGenderFilter, 1) {
		uiState.athletesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Ж##athletesFilter", &uiState.athleteGenderFilter, 2) {
		uiState.athletesDirty = true
	}

	imgui.TextUnformatted("Сортировка")
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Имя##athletesSort", &uiState.athletesSortSwitch, 1) {
		uiState.athletesSortFunc = func(a *Athlete, b *Athlete) bool {
			return a.Name < b.Name
		}
		uiState.athletesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Пол##athletesSort", &uiState.athletesSortSwitch, 2) {
		uiState.athletesSortFunc = func(a *Athlete, b *Athlete) bool {
			return a.Gender < b.Gender
		}
		uiState.athletesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("День рождения##athletesSort", &uiState.athletesSortSwitch, 3) {
		uiState.athletesSortFunc = func(a *Athlete, b *Athlete) bool {
			return a.Birthday.Unix() < b.Birthday.Unix()
		}
		uiState.athletesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Страна##athletesSort", &uiState.athletesSortSwitch, 4) {
		uiState.athletesSortFunc = func(a *Athlete, b *Athlete) bool {
			return a.CountryName < b.CountryName
		}
		uiState.athletesDirty = true
	}

	if uiState.athletesDirty {
		uiState.athletesDirty = false
		processAthletes()
	}

	showTable("##athletesTable", []string{ "", "Имя", "Пол", "День рождения", "Страна" },
		uiState.athletesListProcessed, func(a *Athlete) {
			var gender string
			if a.Gender == "M" {
				gender = "М"
			} else {
				gender = "Ж"
			}

			imgui.TableNextRow()
			imgui.TableNextColumn()
			if (imgui.Button("x")) {
				deleteAthlete(a.ID)
				uiState.athletesList, _ = getAthletes()
				uiState.athletesDirty = true
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.Name)
			imgui.TableNextColumn()
			imgui.TextUnformatted(gender)
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.Birthday.String())
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.CountryName)
		})
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

