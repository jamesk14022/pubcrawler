function containsObject(obj, list) {
  var i;
  for (i = 0; i < list.length; i++) {
    if (list[i] === obj) {
      return true;
    }
  }

  return false;
}

function copy(text) {
  return new Promise((resolve, reject) => {
    // 1) Modern Clipboard API approach
    if (navigator?.clipboard?.writeText) {
      navigator.clipboard
        .writeText(text)
        .then(resolve)
        .catch((err) => {
          console.error(
            "Clipboard API failed, falling back to execCommand.",
            err,
          );
          fallbackCopyToClipboard(text, resolve, reject);
        });
    } else {
      // 2) Fallback for older browsers, iOS Safari, some Android browsers, etc.
      fallbackCopyToClipboard(text, resolve, reject);
    }
  });
}

function fallbackCopyToClipboard(text, resolve, reject) {
  // Create a temporary <textarea> to select and copy from
  const textArea = document.createElement("textarea");
  textArea.value = text;

  // Make it off-screen (and fixed) so iOS doesn’t scroll to it
  textArea.style.position = "fixed";
  textArea.style.top = "-9999px";
  textArea.style.left = "-9999px";

  document.body.appendChild(textArea);

  textArea.focus();
  textArea.select();

  try {
    const successful = document.execCommand("copy");
    document.body.removeChild(textArea);

    if (successful) {
      resolve();
    } else {
      reject(new Error("Fallback: Copy command was unsuccessful"));
    }
  } catch (err) {
    document.body.removeChild(textArea);
    reject(err);
  }
}

function updateURL(location, targetPubs, targetAttractions, ...markers) {
  var state = {
    location: location,
    targetN: targetPubs + targetAttractions,
  };
  for (let i = 0; i < markers.length; i++) {
    state[`marker${i + 1}`] = markers[i];
  }

  var pathState = `?location=${location}&target_pubs=${targetPubs}&target_attractions=${targetAttractions}`;
  for (let i = 0; i < markers.length; i++) {
    pathState += `&marker${i + 1}=${markers[i]}`;
  }

  history.pushState(state, "Route", pathState);
}

function convertToGeoJSON(dataArray) {
  return {
    type: "geojson",
    data: {
      type: "FeatureCollection",
      features: dataArray.map((item) => ({
        type: "Feature",
        properties: {
          name: item.name,
          place_id: item.place_id,
          price_level: item.price_level,
          rating: item.rating,
          types: item.types,
          photos: item.photos,
        },
        geometry: {
          type: "Point",
          coordinates: [item.geometry.location.lng, item.geometry.location.lat],
        },
      })),
    },
  };
}

export { containsObject, copy, updateURL, convertToGeoJSON };
