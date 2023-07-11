import * as React from "react";
import Head from "next/head";

import { Dialog } from "@headlessui/react";

import MapLoadingHolder from "../components/map-loading-holder";
import MapboxMap from "../components/mapbox-map";
import { Modal } from "../components/modal";
import { Button, message, Space, Spin } from "antd";
import { LoadingOutlined } from "@ant-design/icons";
import { APIFunctions } from "./api/axios-client";
import { MapContext } from "../context/context";
import { DataUsers } from "./api/type";

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

const antIcon = <LoadingOutlined style={{ fontSize: 24 }} spin />;

function App() {
  const [loading, setLoading] = React.useState(true);
  const handleMapLoading = () => setLoading(false);
  const [isOpen, setOpen] = React.useState(true);

  const [loadingForm, setLoadingForm] = React.useState(false);
  const [messageApi, contextHolder] = message.useMessage();

  const { userCurrent, setUserCurrent, listUsers, setListUsers } =
    React.useContext(MapContext);

  const handleSubmit = async (event: any) => {
    try {
      event.preventDefault();
      setLoadingForm(true);
      const data = {
        username: event.target.username.value,
        radius: Number(event.target.radius.value),
        latitude: Number(event.target.latitude.value) || null,
        longitude: Number(event.target.longitude.value) || null,
      };

      //If do not enter the coordinates, the default location will be taken
      if (!data.latitude || !data.longitude) {
        if (navigator.geolocation) {
          await new Promise((resolve, reject) => {
            navigator.geolocation.getCurrentPosition(function (position) {
              var lng = !data.longitude
                ? position.coords.longitude
                : data.longitude;
              var lat = !data.latitude
                ? position.coords.latitude
                : data.latitude;
              data.latitude = lat;
              data.longitude = lng;
              resolve(null);
            });
          });
        }
      }

      //Loading form
      const res: any = await APIFunctions.post(
        `${process.env.NEXT_PUBLIC_API_ENDPOINT}/user`,
        data
      );

      if (res?.data?.data) {
        localStorage.setItem("user", JSON.stringify(res.data.data));
      }

      const user = res.data.data as DataUsers;
      setUserCurrent(user);
      setListUsers(fakeData);

      setLoadingForm(false);
      setOpen(false);

      messageApi.open({
        type: "success",
        content: "Register success",
      });
    } catch (error) {
      setLoadingForm(false);
      messageApi.open({
        type: "error",
        content: "Err register",
      });
    }
  };

  //Get current user
  React.useEffect(() => {
    const getInfoUser: string | null = localStorage.getItem("user");
    if (getInfoUser) {
      const user: DataUsers = JSON.parse(getInfoUser);
      setUserCurrent(user);
      setListUsers(fakeData);
      setOpen(false);
    }
  }, []);
  
  return (
    <>
      <Head>
        <title>Map</title>
      </Head>
      <div className="app-container">
        <div className="map-wrapper">
          <MapboxMap
            initialOptions={{ center: [38.0983, 55.7038] }}
            onLoaded={handleMapLoading}
          />
        </div>
        {loading && <MapLoadingHolder className="loading-holder" />}
      </div>
      <Modal isOpen={isOpen}>
        <div className="flex h-full items-center justify-center">
          <Dialog.Panel className="w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-10 text-left align-middle shadow-xl transition-all flex flex-col justify-center items-center">
            <form onSubmit={handleSubmit}>
              <div className="w-[500px] h-[430px] flex flex-col justify-start items-start px-12">
                <div className="w-full flex justify-center items-center text-3xl text-black font-black mb-3">
                  {"EBN TEAM"}
                </div>
                <label className="text-lg font-bold text-black">
                  {"Username (*)"}
                </label>
                <input
                  required
                  type="text"
                  id="username"
                  name="username"
                  placeholder="username"
                  className="px-4 py-3 rounded-lg border border-solid border-black w-full mt-3"
                />
                <label className="mt-5 text-lg font-bold text-black">
                  {"Radius"}
                </label>
                <input
                  required
                  type="text"
                  id="radius"
                  name="radius"
                  placeholder="5 .km"
                  className="px-4 py-3 rounded-lg border border-solid border-black w-full mt-3"
                />
                <label className="mt-5 text-lg font-bold text-black">
                  {"Your location (optional)"}
                </label>
                <div className="w-full flex flex-row mt-3 gap-2">
                  <input
                    id="longitude"
                    name="longitude"
                    placeholder="Longitude"
                    className="w-full px-4 py-3 rounded-lg border border-solid border-black "
                  />
                  <input
                    id="latitude"
                    name="latitude"
                    placeholder="Latitude"
                    className="w-full px-4 py-3 rounded-lg border border-solid border-black "
                  />
                </div>
                <button
                  type="submit"
                  // onClick={() => setOpen(false)}
                  disabled={loadingForm && true}
                  className="py-3 w-full mt-8 rounded-lg bg-black text-white font-bold text-lg hover:bg-gray-600"
                >
                  {"Join Application"}
                  {loadingForm && (
                    <Spin className="mr-2" indicator={antIcon} />
                  )}{" "}
                </button>
              </div>
            </form>
          </Dialog.Panel>
        </div>
      </Modal>
    </>
  );
}

export default App;
