import React, { useContext } from "react";
import { MapContext } from "../context/context";
import mapboxgl from "mapbox-gl";
import { DataUsers } from "../pages/api/type";
import { PopupMap } from "../ultis/popupMap";

interface Props {
  map?: mapboxgl.Map;
  dataUsers: DataUsers[];
}

const MakerMap = ({ map, dataUsers }: Props) => {
  const { lng, lat, setLat, setLng } = useContext(MapContext);

  React.useEffect(() => {
    if (map && lat && lng) {
      let currentLocation = [lng, lat] as [number, number];
      map.setCenter([103.86144, 20.994638]);
      const marker1 = new mapboxgl.Marker()
        .setLngLat(currentLocation)
        .addTo(map);

      /**
       * Create a marker
       *  add it to the map.
       **/
      dataUsers.map((item) => {
        const coordinates = [item.longitude, item.latitude] as [number, number];
        const maker = new mapboxgl.Marker({ offset: [0, -23] })
          .setLngLat(coordinates)
          .setPopup(
            new mapboxgl.Popup({ offset: 25 }) // add popups
              .setHTML(PopupMap)
          )
          .addTo(map);

        maker.getElement().addEventListener("click", (event) => {
          map.flyTo({
            center: coordinates,
            zoom: 11,
          });
          maker.togglePopup();
          const popupElement = maker.getPopup().getElement();

          const closeButton = popupElement?.querySelector(".close-button");
          // Add a click event listener to the close button
          closeButton?.addEventListener("click", () => {
            // Remove the popup from the map
            maker.getPopup().remove();
          });
        });

        
        // Get the popup element
        // Get the close button element inside the popup
      });
    }
  }, [map, lat, lng]);
  return <div></div>;
};

export default MakerMap;
