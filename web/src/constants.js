const TIME_SPENT_BAR = 30;

// token scoped and safe for FE use
const MAPBOX_TOKEN =
  "pk.eyJ1IjoiamFtZXNrMTQwMjIiLCJhIjoiY2x2cnZqZnV5MHdnYTJxcXpkOHUybzdrZCJ9.UVs8BFzWjaZVrz7Gc0_Wpg";

const BASE_URL =
  window.location.host === "127.0.0.1:8080" ||
  window.location.host === "localhost:8080"
    ? "http://localhost:8080"
    : "https://pubcrawler.app";

const GOOGLE_MAP_BASE_URL = "https://www.google.com/maps/search";

const maximumLocationCount = 5;
const INITIAL_LOCATION = [-6.2603, 53.3498];

const container = document.getElementById("container");
const refreshButton = document.getElementById("refresh-button");
const shareButton = document.getElementById("shareButton");
const googleMapsButton = document.getElementById("googleMapsButton");
const searchBox = document.getElementById("search-box");
const modalExitButton = document.getElementById("modal-exit");
const noPubsConent = document.getElementById("no-results");
const rightBar = document.getElementById("rightBar");
const nav = document.getElementById("listing-group");
const dataList = document.getElementById("citiesDropdown");
const dropdownCities = document.getElementById("dropdownCities");

const pubMinus = document.querySelectorAll(".marker-quantity-btn-minus");
const pubPlus = document.querySelectorAll(".marker-quantity-btn-plus");

const attractionMinus = document.querySelectorAll(
  ".attraction-quantity-btn-minus",
);
const attractionPlus = document.querySelectorAll(
  ".attraction-quantity-btn-plus",
);

const filterReset = document.querySelectorAll(".filter-reset");

const attractionCounter = document.querySelectorAll(".num-attractions");
const markerCounter = document.querySelectorAll(".num-markers");

const sidebar = document.getElementById("collap-sidebar");
const sidebarToggle = document.getElementById("sidebarToggle");
const secondaryBar = document.getElementById("secondaryBar");

const selectStart = document.getElementById("pointStart");
const locationCountError = document.getElementById("location-count-error");

export {
  TIME_SPENT_BAR,
  MAPBOX_TOKEN,
  BASE_URL,
  GOOGLE_MAP_BASE_URL,
  INITIAL_LOCATION,
  container,
  refreshButton,
  shareButton,
  googleMapsButton,
  searchBox,
  modalExitButton,
  noPubsConent,
  rightBar,
  nav,
  dataList,
  pubMinus,
  pubPlus,
  attractionMinus,
  attractionPlus,
  filterReset,
  attractionCounter,
  markerCounter,
  sidebar,
  sidebarToggle,
  // closeBtn,
  selectStart,
  secondaryBar,
  locationCountError,
  maximumLocationCount,
  dropdownCities,
};
