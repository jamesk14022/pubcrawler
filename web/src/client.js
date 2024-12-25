import { TIME_SPENT_BAR, maximumLocationCount } from "./constants.js";
import { containsObject, copy, updateURL } from "./utils.js";
import { getCityPoints, postCrawl, getCities, getPubs } from "./api.js";
import {
  clearBarInformationBox,
  clearCityList,
  populateCityList,
  setShareButtonCopied,
  setupShareButtonEvents,
  setupRefreshButtonEvents,
  setupModalExitButtonEvents,
  setRouteDuration,
  setMarkersDisplay,
  setAttractionDisplay,
  showRightBar,
  hideLoading,
  showLoading,
  toggleNoPubsResults,
  setupPubPlusMinusEvents,
  setupAttractionPlusMinusEvents,
  setupFilterResetEvent,
  renderBarInformationBox,
  setupPillClosedEvents,
  hidePill,
  showPill,
  setLocationCountError,
  toggleSidebar,
  setupSidebarToggleEvents,
  hideDropdownCities,
} from "./ui.js";
import {
  flyToLocation,
  renderMapRoute,
  renderAlternativeAttractionMarkers,
  removeAlternativeAttractionMarkers,
  renderRouteMarker,
  setupRenderAlternativeAttractionMarkersPopup,
  removeExistingRoute,
  map,
} from "./map.js";

import "flowbite";
import "./styles.css";

let currentLocation = "dublin";
let currentMarkers = [];
let selectedFirstLocation = "";
let selectedFirstLocationType = "";
let selectedPubs = 3;
let selectedAttractions = 1;
let cityPoints = {};

setupShareButtonEvents(() => {
  copyShareLink();
});

setupPubPlusMinusEvents(
  () => {
    if (selectedPubs === 2) {
      return;
    }
    selectedPubs -= 1;
    setMarkersDisplay(selectedPubs);
  },
  () => {
    if (selectedPubs + selectedAttractions === maximumLocationCount) {
      setLocationCountError();
      return;
    }
    selectedPubs += 1;
    setMarkersDisplay(selectedPubs);
  },
);

setupAttractionPlusMinusEvents(
  () => {
    if (
      selectedFirstLocationType === "attraction" &&
      selectedAttractions === 1
    ) {
      return;
    }
    if (selectedAttractions === 0) {
      return;
    }
    selectedAttractions -= 1;
    setAttractionDisplay(selectedAttractions);
  },
  () => {
    if (selectedPubs + selectedAttractions === maximumLocationCount) {
      setLocationCountError();
      return;
    }
    selectedAttractions += 1;
    setAttractionDisplay(selectedAttractions);
  },
);

setupFilterResetEvent(() => {
  selectedAttractions = 1;
  setAttractionDisplay(selectedAttractions);

  selectedPubs = 3;
  setMarkersDisplay(selectedPubs);
});

export function selectStartEvent(place_id, place_name, type) {
  if (type === "attraction" && selectedAttractions === 0) {
    if (selectedPubs + selectedAttractions === maximumLocationCount) {
      selectedPubs -= 1;
      setMarkersDisplay(selectedPubs);
    }
    selectedAttractions += 1;
    setAttractionDisplay(selectedAttractions);
  }

  if (place_id === "") {
    hidePill();
  } else {
    showPill(place_name);
  }
  selectedFirstLocation = place_id;
  selectedFirstLocationType = type;
}

const clearExistingRoute = () => {
  removeExistingRoute();
  clearBarInformationBox();
  removeAlternativeAttractionMarkers();
  if (currentMarkers !== null) {
    for (let i = currentMarkers.length - 1; i >= 0; i--) {
      currentMarkers[i].remove();
    }
    currentMarkers = [];
  }
};

function copyShareLink() {
  copy(window.location.href);
  setShareButtonCopied();
}

function updateRouteMetrics(e) {
  if (e !== undefined) {
    setRouteDuration(
      parseInt(e[0].duration / 60 + selectedPubs * TIME_SPENT_BAR),
    );
  }
}

async function addAlternativeBarMarkers(route_points) {
  let cityPoints = await getCityPoints(currentLocation);
  cityPoints = await cityPoints.filter(
    (cityPoint) =>
      !containsObject(
        cityPoint.place_id,
        route_points.map((x) => x.place_id),
      ),
  );
  renderAlternativeAttractionMarkers(cityPoints);
  await setupRenderAlternativeAttractionMarkersPopup();
}

async function pageStart() {
  showLoading();
  addCityLocations();

  setAttractionDisplay(selectedAttractions);
  setMarkersDisplay(selectedPubs);

  // Check if the URL contains a query string
  const urlParams = new URLSearchParams(window.location.search);
  if (
    urlParams.has("location") &&
    urlParams.has("target_pubs") &&
    urlParams.has("target_attractions") &&
    urlParams.has("marker1")
  ) {
    setAttractionDisplay(parseFloat(urlParams.get("target_attractions")));
    setMarkersDisplay(parseInt(urlParams.get("target_pubs")));

    // Get the query string values
    const location = urlParams.get("location");
    const targetPubs = parseInt(urlParams.get("target_pubs"));
    const targetAttractions = parseInt(urlParams.get("target_attractions"));
    const markers = [];
    for (let i = 1; i <= targetPubs + targetAttractions; i++) {
      markers.push(urlParams.get(`marker${i}`));
    }

    currentLocation = location.toLowerCase();

    map.on("load", async function () {
      let waypoints = await postCrawl(currentLocation, markers);
      selectedPubs = targetPubs;
      updateRouteMetrics();
      await renderRoute(waypoints);
      hideLoading();
    });
  } else {
    map.on("load", async function () {
      clearExistingRoute();
      showLoading();
      let waypoints = await getPubs(
        selectedPubs,
        selectedAttractions,
        currentLocation,
        selectedFirstLocation,
      );
      await renderRoute(waypoints);
      updateRouteMetrics();
      hideLoading();
    });
  }
}

async function renderRoute(waypoints) {
  clearExistingRoute();

  waypoints.forEach((waypoint, index) => {
    let m = renderRouteMarker(waypoint, index);
    currentMarkers.push(m);
    renderBarInformationBox(waypoint, index);
  });

  await addAlternativeBarMarkers(waypoints);
  showRightBar();
  if (waypoints.length !== 0) {
    renderMapRoute(waypoints);

    updateURL(
      currentLocation,
      selectedPubs,
      selectedAttractions,
      ...waypoints.map((waypoint) => waypoint.place_id),
    );
  } else {
    toggleNoPubsResults();
  }
}

function addCityLocations() {
  clearCityList();
  getCities().then((cities) => {
    cityPoints = cities;
    populateCityList(cityPoints);
  });
  updateRouteMetrics();
}

export async function setCity(e) {
  hideDropdownCities();
  hidePill();
  
  let cityName = e.target.dataset.city;
  currentLocation = cityName;
  selectedFirstLocation = "";

  flyToLocation(cityPoints[cityName]);
  clearExistingRoute();
  showLoading();
  let waypoints = await getPubs(
    selectedPubs,
    selectedAttractions,
    currentLocation,
    selectedFirstLocation,
  );
  await renderRoute(waypoints);
  updateRouteMetrics();

  addCityLocations();
  hideLoading();
}

setupPillClosedEvents(async () => {
  selectStartEvent("", "", "");
  hidePill();
});

setupRefreshButtonEvents(async () => {
  clearExistingRoute();
  showLoading();
  let waypoints = await getPubs(
    selectedPubs,
    selectedAttractions,
    currentLocation,
    selectedFirstLocation,
  );
  await renderRoute(waypoints);
  updateRouteMetrics();
  hideLoading();
});

setupSidebarToggleEvents(async () => {
  toggleSidebar();
});

setupModalExitButtonEvents(async () => {
  toggleNoPubsResults();
});

window.onload = pageStart;
