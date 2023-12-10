import React, { useState, useEffect } from "react";
import BoxList from "./BoxList";
import config from "../config";
import { Navigate } from "react-router-dom";

const ApiFetch = () => {
  const [filenames, setFilenames] = useState([]);
  const [redirect, setRedirect] = useState(false);
  useEffect(() => {
    fetch(`${config.fqdn}/api/list`, { method: "GET", mode: "cors" })
      .then((res) => res.json())
      .then((data) => {
        setFilenames(data["list"]);
        setRedirect(data["list"].length === 0 ? true : false);
      });
  }, []);

  const deleteFile = (value) => {
    setFilenames(filenames.filter((file) => file !== value));
    fetch(`${config.fqdn}/api/deleteFile/${value}`, {
      method: "DELETE",
      mode: "cors",
    })
      .then((res) => res.json())
      .then((data) => console.log(data));
  };

  if (redirect) {
    return <Navigate to="/" />;
  }

  return (
    <div>
      <BoxList filenames={filenames} deleteFunc={deleteFile} />
    </div>
  );
};

export default ApiFetch;
