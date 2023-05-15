import "mapbox-gl/dist/mapbox-gl.css";

import React, { useState, useEffect, useContext } from "react";

import mapboxgl from "mapbox-gl";
import { MapContext } from "../context/context";
import PopupUser from "./maker-map";
import { SmileOutlined, LogoutOutlined, HomeOutlined } from "@ant-design/icons";
import { PopupMap } from "../ultis/popupMap";
import Sidebar from "./Sidebar";
import { DataUsers } from "../pages/api/type";
import MakerMap from "./maker-map";

interface MapboxMapProps {
  initialOptions?: Omit<mapboxgl.MapboxOptions, "container">;
  onCreated?(map: mapboxgl.Map): void;
  onLoaded?(map: mapboxgl.Map): void;
  onRemoved?(): void;
}

const fakeData: DataUsers[] = [
  {
    longitude: 104.86144,
    latitude: 20.994638,
    radius: 1212,
    status: false,
    username: "Hello World A",
  },
  {
    longitude: 103.86144,
    latitude: 20.994638,
    radius: 1212,
    status: false,
    username: "Hello World B ",
  },
  {
    longitude: 102.86144,
    latitude: 20.994638,
    radius: 1212,
    status: false,
    username: "Hello World C",
  },
];

function MapboxMap({
  initialOptions = {},
  onCreated,
  onLoaded,
  onRemoved,
}: MapboxMapProps) {
  const [map, setMap] = React.useState<mapboxgl.Map>();
  const { userCurrent, listUsers } = useContext(MapContext);

  const mapNode = React.useRef(null);

  //Display current location in map
  // React.useEffect(() => {
  //   if (map && lat && lng) {
  //     let currentLocation = [lng, lat] as [number, number];
  //     map.on("load", () => {
  //       map.setCenter([103.86144, 20.994638]);
  //       const marker1 = new mapboxgl.Marker()
  //         .setLngLat(currentLocation)
  //         .addTo(map);

  //       /**
  //        * Create a marker
  //        *  add it to the map.
  //        **/
  //       fakeData.map((item) => {
  //         const coordinates = [item.longitude, item.latitude] as [
  //           number,
  //           number
  //         ];
  //         const maker = new mapboxgl.Marker({ offset: [0, -23] })
  //           .setLngLat(coordinates)
  //           .setPopup(
  //             new mapboxgl.Popup({ offset: 25 }) // add popups
  //               .setHTML(PopupMap)
  //           )
  //           .addTo(map);

  //         maker.getElement().addEventListener("click", (event) => {
  //           map.flyTo({
  //             center: coordinates,
  //             zoom: 11,
  //           });
  //           maker.togglePopup();
  //           const popupElement = maker.getPopup().getElement();

  //           const closeButton = popupElement?.querySelector(".close-button");
  //           // Add a click event listener to the close button
  //           closeButton?.addEventListener("click", () => {
  //             // Remove the popup from the map
  //             maker.getPopup().remove();
  //           });
  //         });

  //         // Get the popup element

  //         // Get the close button element inside the popup
  //       });
  //     });
  //   }
  // }, [map]);
  //Mount map first
  React.useEffect(() => {
    const node = mapNode.current;
    if (typeof window === "undefined" || node === null) return;
    let mapboxMap = new mapboxgl.Map({
      container: node,
      accessToken: process.env.NEXT_PUBLIC_MAPBOX_TOKEN,
      style: "mapbox://styles/mapbox/streets-v11",
      // initial position in [lon, lat] format
      //Default Ha Noi
      center: [
        userCurrent ? userCurrent.longitude : 105.44,
        userCurrent ? userCurrent.latitude : 20.53,
      ],
      // initial zoom
      zoom: 14,
      ...initialOptions,
    });
    // Add geolocate control to the map.
    mapboxMap.addControl(
      new mapboxgl.GeolocateControl({
        positionOptions: {
          enableHighAccuracy: true,
        },
        // When active the map will receive updates to the device's location as it changes.
        trackUserLocation: true,
        // Draw an arrow next to the location dot to indicate which direction the device is heading.
        showUserHeading: true,
      })
    );

    setMap(mapboxMap);
    if (onCreated) onCreated(mapboxMap);
    if (onLoaded) mapboxMap.once("load", () => onLoaded(mapboxMap));

    return () => {
      mapboxMap.remove();
      setMap(undefined);
      if (onRemoved) onRemoved();
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div ref={mapNode} style={{ width: "100%", height: "100%" }}>
      <MakerMap map={map} dataUsers={fakeData} />
      <Sidebar />
    </div>
  );
}

export default MapboxMap;
