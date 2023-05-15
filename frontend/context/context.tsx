import { Context, createContext, useState } from "react";
import mapboxgl from "mapbox-gl";
import { DataUsers } from "../pages/api/type";

interface MapboxMapProps {
  initialOptions?: Omit<mapboxgl.MapboxOptions, "container">;
  onCreated?(map: mapboxgl.Map): void;
  onLoaded?(map: mapboxgl.Map): void;
  onRemoved?(): void;
}

interface MapboxMap {
  map: mapboxgl.Map | undefined;
  lng: number;
  lat: number;
  listUsers: DataUsers[] | [];
  userCurrent: DataUsers | undefined;
  setLng: (lng: number) => void;
  setLat: (lat: number) => void;
  setMap: (map: mapboxgl.Map | undefined) => void;
  setListUsers: (users: DataUsers[]) => void;
  setUserCurrent: (user: DataUsers) => void;
}

export const MapContext: Context<MapboxMap> = createContext<MapboxMap>({
  lng: 0,
  lat: 0,
  map: undefined,
  listUsers: [],
  userCurrent: undefined,
  setLng: () => {},
  setLat: () => {},
  setMap: () => {},
  setListUsers: () => {},
  setUserCurrent: () => {},
});

export const MapContextProvider = (props: any) => {
  const [lng, setLng] = useState<number>(0);
  const [lat, setLat] = useState<number>(0);
  const [map, setMap] = useState<mapboxgl.Map>();
  const [listUsers, setListUsers] = useState<DataUsers[]>([]);
  const [userCurrent, setUserCurrent] = useState<DataUsers>();
  return (
    <MapContext.Provider
      value={{
        lng: lng,
        setLng: setLng,
        lat: lat,
        setLat: setLat,
        map: map,
        setMap: setMap,
        listUsers: listUsers,
        setListUsers: setListUsers,
        userCurrent: userCurrent,
        setUserCurrent: setUserCurrent,
      }}
    >
      {props.children}
    </MapContext.Provider>
  );
};
