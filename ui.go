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
	TabSports
	TabAthletes
	TabSites
	TabTeams
	TabCompetitions
	TabMedals
)

var tabs []Tab = []Tab{TabCountries, TabSports, TabAthletes, TabSites, TabTeams, TabCompetitions, TabMedals}

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

	sitesDirty bool
	sitesList []Site
	sitesListProcessed []*Site
	siteNameInput string
	siteNameFilter string

	teamsDirty bool
	teamsList []Team
	teamsListProcessed []*Team
	teamNameInput string
	teamNameFilter string
	teamCountryInput Country
	teamSportInput Sport
	teamMemberSelection map[int]int

	competitionsDirty bool
	competitionsList []Competition
	competitionsListProcessed []*Competition
	competitionDateInput string
	competitionTimeInput string
	competitionSportInput Sport
	competitionSiteInput Site
	competitionFilterSport Sport
	competitionFilterSite  Site

    medalsList []CountryMedals
    medalsListProcessed []*CountryMedals
    medalsDirty bool
    medalsSortSwitch int32
    medalsSortFunc func(a, b *CountryMedals) bool

	hasError bool
	error string
}

var tableFlags imgui.TableFlags =
	imgui.TableFlagsBordersOuter | imgui.TableFlagsBordersInner | imgui.TableFlagsRowBg |
		imgui.TableFlagsScrollY | imgui.TableFlagsScrollX | imgui.TableFlagsResizable |
		imgui.TableFlagsHighlightHoveredColumn

func (tab Tab) name() string {
	switch tab {
	case TabCountries: return "Страны"
	case TabSports: return "Виды спорта"
	case TabAthletes: return "Спортсмены"
	case TabSites: return "Места проведения"
	case TabTeams: return "Команды"
	case TabCompetitions: return "Соревнования"
	case TabMedals: return "Медали"
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
	case TabSites: showSites(switched)
	case TabTeams: showTeams(switched)
	case TabCompetitions: showCompetitions(switched)
	case TabMedals: showMedals(switched)
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

func processCompetitions() {
	uiState.competitionsListProcessed = make([]*Competition, 0, len(uiState.competitionsList))
	for i := range uiState.competitionsList {
		c := &uiState.competitionsList[i]
		if uiState.competitionFilterSport.Code != "" &&
			uiState.competitionFilterSport.Code != c.Sport.Code {
			continue
		}
		if uiState.competitionFilterSite.ID != 0 &&
			uiState.competitionFilterSite.ID != c.Site.ID {
			continue
		}
		uiState.competitionsListProcessed = append(uiState.competitionsListProcessed, c)
	}
}

func showCompetitions(switched bool) {
	if switched {
		uiState.competitionsList, _ = getCompetitions()
		uiState.competitionsDirty = true
	}

	avail := imgui.ContentRegionAvail()

	imgui.SetNextItemWidth(avail.X / 5)
	imgui.InputTextWithHint("##compDate", "Дата (YYYY-MM-DD)", &uiState.competitionDateInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 5)
	imgui.InputTextWithHint("##compTime", "Время (HH:MM)", &uiState.competitionTimeInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 5)
	pickSport(&uiState.competitionSportInput, "##pickSportCombo")
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 5)
	pickSite(&uiState.competitionSiteInput, "##pickSiteCombo")
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 5)
	if imgui.Button("Добавить") {

		if uiState.competitionDateInput == "" || uiState.competitionTimeInput == "" {
			showError(fmt.Errorf("Введите дату и время"))
		} else {
			t, err := time.Parse("2006-01-02 15:04",
				uiState.competitionDateInput+" "+uiState.competitionTimeInput)

			if err != nil {
				showError(err)
			} else {
				err = addCompetition(t, uiState.competitionSportInput.Code, uiState.competitionSiteInput.ID)
				if err != nil {
					showError(err)
				} else {
					uiState.competitionsList, _ = getCompetitions()
					uiState.competitionsDirty = true
				}
			}
		}
	}

	imgui.Separator()

	imgui.Text("Фильтр")

	imgui.Text("Спорт")
	imgui.SetNextItemWidth(avail.X / 3)
	imgui.SameLine()
	if pickSport(&uiState.competitionFilterSport, "##pickSportComboFilter") {
		uiState.competitionsDirty = true
	}
	imgui.SameLine()
	if imgui.Button("x##clearSportFilter") {
		uiState.competitionFilterSport = Sport{}
		uiState.competitionsDirty = true
	}

	imgui.Text("Место")
	imgui.SetNextItemWidth(avail.X / 3)
	imgui.SameLine()
	if pickSite(&uiState.competitionFilterSite, "##pickSiteComboFilter") {
		uiState.competitionsDirty = true
	}
	imgui.SameLine()
	if imgui.Button("x##clearSiteFilter") {
		uiState.competitionFilterSite = Site{}
		uiState.competitionsDirty = true
	}

	if uiState.competitionsDirty {
		uiState.competitionsDirty = false
		processCompetitions()
	}

	showTable("##competitionsTable",
		[]string{"", "Дата и время", "Вид спорта", "Место"},
		uiState.competitionsListProcessed,
		func(c *Competition) {

			imgui.TableNextRow()

			imgui.TableNextColumn()
			if imgui.Button("x") {
				deleteCompetition(c.ID)
				uiState.competitionsList, _ = getCompetitions()
				uiState.competitionsDirty = true
			}

			imgui.TableNextColumn()
			imgui.TextUnformatted(c.Time.Format("2006-01-02 15:04"))

			imgui.TableNextColumn()
			imgui.TextUnformatted(c.Sport.Name)

			imgui.TableNextColumn()
			imgui.TextUnformatted(c.Site.Name)
		})
}

func processSites() {
	uiState.sitesListProcessed = make([]*Site, 0, len(uiState.sitesList))

	for i := range uiState.sitesList {
		s := &uiState.sitesList[i]
		if len(uiState.siteNameFilter) > 0 && !strings.Contains(s.Name, uiState.siteNameFilter) {
			continue
		}
		uiState.sitesListProcessed = append(uiState.sitesListProcessed, s)
	}
}

func showSites(switched bool) {
	if switched {
		uiState.sitesList, _ = getSites()
		uiState.sitesDirty = true
	}

	avail := imgui.ContentRegionAvail()
	imgui.SetNextItemWidth(avail.X / 3)
	imgui.InputTextWithHint("##siteNameInput", "Название места", &uiState.siteNameInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 1)
	if imgui.Button("Добавить") {
		err := addSite(uiState.siteNameInput)
		if err != nil {
			showError(err)
		} else {
			uiState.sitesList, _ = getSites()
			uiState.sitesDirty = true
		}
	}

	imgui.Separator()
	imgui.TextUnformatted("Фильтр")

	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 3)
	if imgui.InputTextWithHint("##siteNameFilter", "Название", &uiState.siteNameFilter, 0, nil) {
		uiState.sitesDirty = true
	}

	if uiState.sitesDirty {
		uiState.sitesDirty = false
		processSites()
	}

	showTable("##sitesTable", []string{"", "Название"},
		uiState.sitesListProcessed, func(s *Site) {
			imgui.TableNextRow()
			imgui.TableNextColumn()
			if imgui.Button("x") {
				deleteSite(s.ID)
				uiState.sitesList, _ = getSites()
				uiState.sitesDirty = true
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(s.Name)
		})
}

func processAthletes() {
	uiState.athletesListProcessed = make([]*Athlete, 0, len(uiState.athletesList))

	for i := range uiState.athletesList {
		a := &uiState.athletesList[i]
		if len(uiState.athleteNameFilter) > 0 && !strings.Contains(a.Name, uiState.athleteNameFilter) {
			continue
		} else if uiState.athleteGenderFilter > 0 {
			if (uiState.athleteGenderFilter == 1 && a.Gender != "M") ||
				(uiState.athleteGenderFilter == 2 && a.Gender != "F") {
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

func pickCountry(country *Country) {
	if imgui.BeginCombo("##pickCountryCombo", country.Name) {
		defer imgui.EndCombo()
		countries, _ := getCountries()
		for _, c := range countries {
			if imgui.SelectableBool(c.Name) {
				*country = c
				return
			}
		}
	}
}

func pickSite(site *Site, id string) bool {
	if imgui.BeginCombo(id, site.Name) {
		defer imgui.EndCombo()
		sites, _ := getSites()
		for _, s := range sites {
			if (imgui.SelectableBool(s.Name)) {
				*site = s
				return true
			}
		}
	}
	return false
}

func pickSport(sport *Sport, id string) bool {
	if imgui.BeginCombo(id, sport.Name) {
		defer imgui.EndCombo()
		sports, _ := getSports()
		for _, s := range sports {
			if (imgui.SelectableBool(s.Name)) {
				*sport = s
				return true
			}
		}
	}
	return false
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
	pickCountry(&uiState.athleteCountryInput)
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

	showTable("##athletesTable", []string{"", "Имя", "Пол", "День рождения", "Страна"},
		uiState.athletesListProcessed, func(a *Athlete) {
			var gender string
			if a.Gender == "M" {
				gender = "М"
			} else {
				gender = "Ж"
			}

			imgui.TableNextRow()
			imgui.TableNextColumn()
			if imgui.Button("x") {
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

func processTeams() {
	uiState.teamsListProcessed = make([]*Team, 0, len(uiState.teamsList))

	for i := range uiState.teamsList {
		t := &uiState.teamsList[i]
		if len(uiState.teamNameFilter) > 0 && !strings.Contains(t.Name, uiState.teamNameFilter) {
			continue
		}
		uiState.teamsListProcessed = append(uiState.teamsListProcessed, t)
	}
}

func processMedals() {
    uiState.medalsListProcessed = make([]*CountryMedals, len(uiState.medalsList))
    for i := range uiState.medalsList {
        uiState.medalsListProcessed[i] = &uiState.medalsList[i]
    }

    if uiState.medalsSortFunc != nil {
        sort.Slice(uiState.medalsListProcessed, func(i, j int) bool {
            return uiState.medalsSortFunc(uiState.medalsListProcessed[i], uiState.medalsListProcessed[j])
        })
    }
}

func showMedals(switched bool) {
    if switched {
        medals, err := getCountryMedals()
        if err != nil {
            showError(err)
        } else {
            uiState.medalsList = medals
        }
        uiState.medalsDirty = true
    }

    imgui.TextUnformatted("Сортировка")

    if imgui.SameLine(); imgui.RadioButtonIntPtr("Серебро##medalsSort", &uiState.medalsSortSwitch, 1) {
        uiState.medalsSortFunc = func(a, b *CountryMedals) bool {
            if a.Gold != b.Gold {
                return a.Gold > b.Gold
            }
            if a.Silver != b.Silver {
                return a.Silver > b.Silver
            }
            return a.Bronze > b.Bronze
        }
        uiState.medalsDirty = true
    }
    if imgui.SameLine(); imgui.RadioButtonIntPtr("По серебру##medalsSort", &uiState.medalsSortSwitch, 2) {
        uiState.medalsSortFunc = func(a, b *CountryMedals) bool {
            if a.Silver != b.Silver {
                return a.Silver > b.Silver
            }
            if a.Gold != b.Gold {
                return a.Gold > b.Gold
            }
            return a.Bronze > b.Bronze
        }
        uiState.medalsDirty = true
    }
    if imgui.SameLine(); imgui.RadioButtonIntPtr("Бронза##medalsSort", &uiState.medalsSortSwitch, 3) {
        uiState.medalsSortFunc = func(a, b *CountryMedals) bool {
            if a.Bronze != b.Bronze {
                return a.Bronze > b.Bronze
            }
            if a.Gold != b.Gold {
                return a.Gold > b.Gold
            }
            return a.Silver > b.Silver
        }
        uiState.medalsDirty = true
    }
    if imgui.SameLine(); imgui.RadioButtonIntPtr("По общему числу##medalsSort", &uiState.medalsSortSwitch, 4) {
        uiState.medalsSortFunc = func(a, b *CountryMedals) bool {
            if a.Total != b.Total {
                return a.Total > b.Total
            }
            if a.Gold != b.Gold {
                return a.Gold > b.Gold
            }
            return a.Silver > b.Silver
        }
        uiState.medalsDirty = true
    }
    if imgui.SameLine(); imgui.RadioButtonIntPtr("Страна##medalsSort", &uiState.medalsSortSwitch, 5) {
        uiState.medalsSortFunc = func(a, b *CountryMedals) bool {
            return a.Country < b.Country
        }
        uiState.medalsDirty = true
    }


    if uiState.medalsDirty {
        uiState.medalsDirty = false
        processMedals()
    }

    imgui.Separator()


    showTable("##medalsTable", []string{"Страна", "Золото", "Серебро", "Бронза", "Всего"},
        uiState.medalsListProcessed, func(m *CountryMedals) {
            imgui.TableNextRow()
            imgui.TableNextColumn()
            imgui.TextUnformatted(m.Country)
            imgui.TableNextColumn()
            imgui.TextUnformatted(fmt.Sprintf("%d", m.Gold))
            imgui.TableNextColumn()
            imgui.TextUnformatted(fmt.Sprintf("%d", m.Silver))
            imgui.TableNextColumn()
            imgui.TextUnformatted(fmt.Sprintf("%d", m.Bronze))
            imgui.TableNextColumn()
            imgui.TextUnformatted(fmt.Sprintf("%d", m.Total))
        })
}
func showTeams(switched bool) {
	if switched {
		uiState.teamsList, _ = getTeams()
		uiState.teamsDirty = true
		if uiState.teamMemberSelection == nil {
			uiState.teamMemberSelection = make(map[int]int)
		}
	}

	avail := imgui.ContentRegionAvail()
	imgui.SetNextItemWidth(avail.X / 4)
	imgui.InputTextWithHint("##teamNameInput", "Название команды", &uiState.teamNameInput, 0, nil)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	pickCountry(&uiState.teamCountryInput)
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	pickSport(&uiState.teamSportInput, "##pickSportCombo")
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 1)
	if imgui.Button("Добавить") {
		err := addTeam(uiState.teamNameInput, uiState.teamCountryInput.Code, uiState.teamSportInput.Code)
		if err != nil {
			showError(err)
		} else {
			uiState.teamsList, _ = getTeams()
			uiState.teamsDirty = true
		}
	}

	imgui.Separator()
	imgui.TextUnformatted("Фильтр")
	imgui.SameLine()
	imgui.SetNextItemWidth(avail.X / 4)
	if imgui.InputTextWithHint("##teamNameFilter", "Название", &uiState.teamNameFilter, 0, nil) {
		uiState.teamsDirty = true
	}

	if uiState.teamsDirty {
		uiState.teamsDirty = false
		processTeams()
	}

	allAthletes, _ := getAthletes()

	showTable("##teamsTable", []string{"", "Название", "Страна", "Вид спорта", "Участники"},
		uiState.teamsListProcessed, func(t *Team) {
			imgui.TableNextRow()
			imgui.TableNextColumn()
			if imgui.Button("x") {
				deleteTeam(t.ID)
				uiState.teamsList, _ = getTeams()
				uiState.teamsDirty = true
			}
			imgui.TableNextColumn()
			imgui.TextUnformatted(t.Name)
			imgui.TableNextColumn()
			imgui.TextUnformatted(t.Country.Name)
			imgui.TableNextColumn()
			imgui.TextUnformatted(t.Sport.Name)
			imgui.TableNextColumn()

			label := fmt.Sprintf("members_%d", t.ID)
			countLabel := fmt.Sprintf("Участники (%d)", len(t.Members))
			if imgui.TreeNodeStr(fmt.Sprintf("%s##%s", countLabel, label)) {
				for i := range t.Members {
					m := t.Members[i]
					if imgui.Button(fmt.Sprintf("X##%d_%d", t.ID, m.ID)) {
						if err := deleteAthleteFromTeam(t.ID, m.ID); err != nil {
							showError(err)
						} else {
							uiState.teamsList, _ = getTeams()
							uiState.teamsDirty = true
						}
					}
					imgui.SameLine()
					imgui.TextUnformatted(m.Name)
				}
				if imgui.Button(fmt.Sprintf("Добавить##add_%d", t.ID)) {
					sel := uiState.teamMemberSelection[t.ID]
					if sel == 0 {
						showError(fmt.Errorf("Не выбран спортсмен"))
					} else {
						if err := addAthleteToTeam(t.ID, sel); err != nil {
							showError(err)
						} else {
							uiState.teamsList, _ = getTeams()
							uiState.teamsDirty = true
						}
					}
				}
				imgui.SameLine()
				comboLabel := fmt.Sprintf("##addMemberCombo%d", t.ID)
				selectedID := uiState.teamMemberSelection[t.ID]
				selectedName := "Выберите спортсмена"
				for _, a := range allAthletes {
					if a.ID == selectedID {
						selectedName = a.Name
						break
					}
				}
				if imgui.BeginCombo(comboLabel, selectedName) {
					for _, a := range allAthletes {
						if imgui.SelectableBool(a.Name) {
							uiState.teamMemberSelection[t.ID] = a.ID
						}
					}
					imgui.EndCombo()
				}
				imgui.TreePop()
			}
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
		if imgui.Button("OK") {
			imgui.CloseCurrentPopup()
		}
		imgui.EndPopup()
	}

	imgui.End()
}

func initUI() {
	uiState.oldTab = 100500
	uiState.teamMemberSelection = make(map[int]int)
}

