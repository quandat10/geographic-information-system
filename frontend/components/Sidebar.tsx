import classNames from "classnames";
import Link from "next/link";
import { useRouter } from "next/router";
import React, { useState, useMemo, FunctionComponent } from "react";
import { Button, Modal } from "antd";
import {
  UserOutlined,
  LogoutOutlined,
  PushpinOutlined,
  CloseOutlined,
} from "@ant-design/icons";
import { DataUsers } from "../pages/api/type";
// const menuItems = [
//   { id: 1, label: "Home", icon: HomeIcon, link: "/" },
//   { id: 2, label: "Manage Posts", icon: ArticleIcon, link: "/posts" },
//   { id: 3, label: "Manage Users", icon: UsersIcon, link: "/users" },
//   { id: 4, label: "Manage Tutorials", icon: VideosIcon, link: "/tutorials" },
// ];

interface Props {
  dataUsers?: DataUsers[];
  map?: mapboxgl.Map;
}

const Sidebar = () => {
  const [toggleCollapse, setToggleCollapse] = useState(false);
  const [isCollapsible, setIsCollapsible] = useState(false);
  const router = useRouter();

  // const activeMenu = useMemo(
  //   () => menuItems.find((menu) => menu.link === router.pathname),
  //   [router.pathname]
  // );

  const wrapperClasses = classNames(
    "h-screen px-4 pt-8 pb-4 bg-black flex justify-between flex-col fixed z-10",
    {
      ["w-80"]: !toggleCollapse,
      ["w-20"]: toggleCollapse,
    }
  );

  const collapseIconClasses = classNames(
    "p-4 rounded bg-light-lighter absolute right-0",
    {
      "rotate-180": toggleCollapse,
    }
  );

  // const getNavItemClasses = (menu) => {
  //   return classNames(
  //     "flex items-center cursor-pointer hover:bg-light-lighter rounded w-full overflow-hidden whitespace-nowrap",
  //     {
  //       ["bg-light-lighter"]: activeMenu.id === menu.id,
  //     }
  //   );
  // };

  const onMouseOver = () => {
    setIsCollapsible(!isCollapsible);
  };

  const handleSidebarToggle = () => {
    setToggleCollapse(!toggleCollapse);
  };

  return (
    <div className="sidebar">
      <div className="heading flex gap-4 flex-col items-center">
        {/* <div className="flex items-center">
          <img className="img-icon" src="avatar-svgrepo-com.svg" alt="" />
        </div> */}
        <h1 className="mt-5 text-base">WELL COME: ABC</h1>
        <div className="flex gap-4 items-center">
          <img className="img-icon" src="avatar-svgrepo-com.svg" alt="" />
          <div className="flex flex-col ">
            <span className="text-indigo-50 my-0">Tom Cook</span>
            <span className="text-indigo-50 my-0">Tom Cook</span>
          </div>
          <Button className="ml-4 " type="primary" ghost>
            Open Modal
          </Button>
        </div>
      </div>
      <div id="listings" className="listings">
        <div id="listing-0" className="item flex items-center">
          <div className="w-11/12">
            <a href="#" className="flex gap-4 items-center title" id="link-0">
              <UserOutlined className="icon-antd pr-5" />
              <p>AGSDAGAG</p>
            </a>
            <div className="flex gap-4 items-center">
              <PushpinOutlined className="icon-antd pr-5" />
              <div>
                <p className="text-sm">Lgo: 23423213123 -------</p>
                <p className="text-sm">Lgo: 23423213123 -----</p>
              </div>
            </div>
          </div>
          <div>
            <CloseOutlined className="icon-antd icon-cancle" />
          </div>
        </div>
        <div id="listing-0" className="item">
          <a href="#" className="title" id="link-0">
            P St NW
          </a>
          <div> 234-7336</div>
        </div>
      </div>
      <div className="footer flex items-center">
        <LogoutOutlined className="icon-antd w-2/6"/>
        <h1 className="text-base">LOG OUT</h1>
      </div>
    </div>
  );
};

export default Sidebar;
