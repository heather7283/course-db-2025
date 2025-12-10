package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

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
	countriesDirty bool
	countriesSortSwitch int32
	countriesSortFunc func(a *Country, b *Country) bool
	countriesListProcessed []*Country
	countryCodeFilter string
	countryNameFilter string

	sportsList []Sport
	sportCodeInput string
	sportNameInput string
	sportIsTeamInput bool
	sportsDirty bool
	sportsSortSwitch int32
	sportsSortFunc func(a *Sport, b *Sport) bool
	sportsListProcessed []*Sport
	sportCodeFilter string
	sportNameFilter string
	sportTeamFilter int32

	athletesDirty bool
	athletesSortSwitch int32
	athletesSortFunc func(a *Athlete, b *Athlete) bool
	athletesList []Athlete
	athletesListProcessed []*Athlete
	athleteNameInput string
	athleteIsMaleInput bool
	athleteBirthdayInput string
	athleteCountryInput Country
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

	avail := imgui.ContentRegionAvail()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##athleteNameInput", "Имя", &uiState.athleteNameInput, 0, nil)
	if imgui.SameLine(); imgui.RadioButtonBool("М", uiState.athleteIsMaleInput) {
		uiState.athleteIsMaleInput = true
	}
	if imgui.SameLine(); imgui.RadioButtonBool("Ж", !uiState.athleteIsMaleInput) {
		uiState.athleteIsMaleInput = false
	}
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.SameLine()
	imgui.InputTextWithHint("##athleteBirthdayInput", "День рождения", &uiState.athleteBirthdayInput, 0, nil)
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.SameLine()
	if imgui.BeginCombo("##athleteCountryInput", uiState.athleteCountryInput.Name) {
		countries, _ := getCountries()
		for _, c := range countries {
			if (imgui.SelectableBool(c.Name)) {
				uiState.athleteCountryInput = c
				break
			}
		}
		imgui.EndCombo()
	}
	imgui.SetNextItemWidth(avail.X / 1)
	imgui.SameLine()
	if imgui.Button("Добавить") {
		if t, err := time.Parse(time.DateOnly, uiState.athleteBirthdayInput); err != nil {
			showError(err)
		} else {
			err := addAthlete(uiState.athleteNameInput, uiState.athleteIsMaleInput, t, uiState.athleteCountryInput.Code)
			if err != nil {
				showError(err)
			} else {
				uiState.athletesList, _ = getAthletes()
				uiState.athletesDirty = true
			}
		}
	}

	imgui.Separator()

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
			imgui.TextUnformatted(a.Birthday.Format(time.DateOnly))
			imgui.TableNextColumn()
			imgui.TextUnformatted(a.CountryName)
		})
}

func processSports() {
	uiState.sportsListProcessed = make([]*Sport, 0, len(uiState.sportsList))

	for i := range uiState.sportsList {
		s := &uiState.sportsList[i]
		if len(uiState.sportCodeFilter) > 0 && !strings.Contains(s.Code, uiState.sportCodeFilter) {
			continue
		}
		if len(uiState.sportNameFilter) > 0 && !strings.Contains(s.Name, uiState.sportNameFilter) {
			continue
		}
		if uiState.sportTeamFilter > 0 {
			if (uiState.sportTeamFilter == 1 && !s.IsTeam) ||
				(uiState.sportTeamFilter == 2 && s.IsTeam) {
				continue
			}
		}
		uiState.sportsListProcessed = append(uiState.sportsListProcessed, s)
	}

	if uiState.sportsSortFunc != nil {
		sort.Slice(uiState.sportsListProcessed, func(i, j int) bool {
			return uiState.sportsSortFunc(uiState.sportsListProcessed[i], uiState.sportsListProcessed[j])
		})
	}
}

func showSports(switched bool) {
	if switched {
		uiState.sportsList, _ = getSports()
		uiState.sportsDirty = true
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
		err := addSport(uiState.sportCodeInput, uiState.sportNameInput, uiState.sportIsTeamInput)
		if err != nil {
			showError(err)
		} else {
			uiState.sportsList, _ = getSports()
			uiState.sportsDirty = true
		}
	}

	imgui.Separator()
	imgui.TextUnformatted("Фильтр")

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##sportCodeFilter", "Код", &uiState.sportCodeFilter, 0, nil) {
		uiState.sportsDirty = true
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##sportNameFilter", "Название", &uiState.sportNameFilter, 0, nil) {
		uiState.sportsDirty = true
	}

	imgui.SameLine()
	if imgui.RadioButtonIntPtr("Все##sportTeamFilter", &uiState.sportTeamFilter, 0) {
		uiState.sportsDirty = true
	}
	imgui.SameLine()
	if imgui.RadioButtonIntPtr("Командные##sportTeamFilter", &uiState.sportTeamFilter, 1) {
		uiState.sportsDirty = true
	}
	imgui.SameLine()
	if imgui.RadioButtonIntPtr("Одиночные##sportTeamFilter", &uiState.sportTeamFilter, 2) {
		uiState.sportsDirty = true
	}

	imgui.TextUnformatted("Сортировка")
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Код##sportsSort", &uiState.sportsSortSwitch, 1) {
		uiState.sportsSortFunc = func(a *Sport, b *Sport) bool {
			return a.Code < b.Code
		}
		uiState.sportsDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Название##sportsSort", &uiState.sportsSortSwitch, 2) {
		uiState.sportsSortFunc = func(a *Sport, b *Sport) bool {
			return a.Name < b.Name
		}
		uiState.sportsDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Тип##sportsSort", &uiState.sportsSortSwitch, 3) {
		uiState.sportsSortFunc = func(a *Sport, b *Sport) bool {
			if a.IsTeam && !b.IsTeam {
				return true
			} else if !a.IsTeam && b.IsTeam {
				return false
			}
			return a.Name < b.Name
		}
		uiState.sportsDirty = true
	}

	if uiState.sportsDirty {
		uiState.sportsDirty = false
		processSports()
	}

	showTable("##sportsTable", []string{"", "Код", "Название", "Тип"},
		uiState.sportsListProcessed, func(s *Sport) {
			imgui.TableNextRow()
			imgui.TableNextColumn()
			if imgui.Button("x") {
				deleteSport(s.Code)
				uiState.sportsList, _ = getSports()
				uiState.sportsDirty = true
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(s.Code)
			imgui.TableNextColumn()
			imgui.TextUnformatted(s.Name)
			imgui.TableNextColumn()
			if s.IsTeam {
				imgui.TextUnformatted("Командный")
			} else {
				imgui.TextUnformatted("Одиночный")
			}
		})
}

func processCountries() {
	uiState.countriesListProcessed = make([]*Country, 0, len(uiState.countriesList))

	for i := range uiState.countriesList {
		c := &uiState.countriesList[i]
		if len(uiState.countryCodeFilter) > 0 && !strings.Contains(c.Code, uiState.countryCodeFilter) {
			continue
		}
		if len(uiState.countryNameFilter) > 0 && !strings.Contains(c.Name, uiState.countryNameFilter) {
			continue
		}
		uiState.countriesListProcessed = append(uiState.countriesListProcessed, c)
	}

	if uiState.countriesSortFunc != nil {
		sort.Slice(uiState.countriesListProcessed, func(i, j int) bool {
			return uiState.countriesSortFunc(uiState.countriesListProcessed[i], uiState.countriesListProcessed[j])
		})
	}
}

func showCountries(switched bool) {
	if switched {
		uiState.countriesList, _ = getCountries()
		uiState.countriesDirty = true
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
			uiState.countriesDirty = true
		}
	}

	imgui.Separator()
	imgui.TextUnformatted("Фильтр")

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##countryCodeFilter", "Код", &uiState.countryCodeFilter, 0, nil) {
		uiState.countriesDirty = true
	}

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##countryNameFilter", "Название", &uiState.countryNameFilter, 0, nil) {
		uiState.countriesDirty = true
	}

	imgui.TextUnformatted("Сортировка")
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Код##countriesSort", &uiState.countriesSortSwitch, 1) {
		uiState.countriesSortFunc = func(a *Country, b *Country) bool {
			return a.Code < b.Code
		}
		uiState.countriesDirty = true
	}
	if imgui.SameLine(); imgui.RadioButtonIntPtr("Название##countriesSort", &uiState.countriesSortSwitch, 2) {
		uiState.countriesSortFunc = func(a *Country, b *Country) bool {
			return a.Name < b.Name
		}
		uiState.countriesDirty = true
	}

	if uiState.countriesDirty {
		uiState.countriesDirty = false
		processCountries()
	}

	showTable("##countriesTable", []string{"", "Код", "Название"},
		uiState.countriesListProcessed, func(c *Country) {
			imgui.TableNextRow()
			imgui.TableNextColumn()
			if imgui.Button("x") {
				deleteCountry(c.Code)
				uiState.countriesList, _ = getCountries()
				uiState.countriesDirty = true
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(c.Code)
			imgui.TableNextColumn()
			imgui.TextUnformatted(c.Name)
		})
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

