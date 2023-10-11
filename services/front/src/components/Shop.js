import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import DeleteIcon from "@mui/icons-material/DeleteOutlined";
import Button from "@mui/material/Button";
import AddIcon from "@mui/icons-material/Add";
import { v4 as uuidv4 } from "uuid";
import config from "../config";

const EditToolbar = (props) => {
  const { setRows } = props;

  const handleClick = () => {
    setRows((oldRows) => [
      {
        randomId: uuidv4(),
        id: "",
        site_name: "",
        name: "",
        interval: "",
        url: "",
      },
      ...oldRows,
    ]);
  };

  return (
    <GridToolbarContainer>
      <Button coloer="primary" startIcon={<AddIcon />} onClick={handleClick}>
        Add record
      </Button>
    </GridToolbarContainer>
  );
};

export const Shops = () => {
  const [rows, setRows] = useState([]);

  const handleDeleteClick = (id) => () => {
    setRows(rows.filter((row) => row.id !== id));
  };

  const columns = [
    { field: "id", headerName: "ID", width: 130, editable: true },
    { field: "site_name", headerName: "Site name", width: 130, editable: true },
    { field: "name", headerName: "Shop name", width: 130, editable: true },
    {
      field: "interval",
      headerName: "Interval",
      width: 90,
      editable: true,
      type: "singleSelect",
      valueOptions: ["daily", "weekly"],
    },
    {
      field: "url",
      headerName: "URL",
      width: 300,
      editable: true,
    },
    {
      field: "actions",
      type: "actions",
      headerName: "Actions",
      width: 100,
      cellClassName: "actions",
      getActions: ({ id }) => {
        return [
          <GridActionsCellItem
            icon={<DeleteIcon />}
            label="Delete"
            onClick={handleDeleteClick(id)}
            color="inherit"
          />,
        ];
      },
    },
  ];

  useEffect(() => {
    fetch(`${config.fqdn}/api/shops`, { method: "GET", mode: "cors" })
      .then((res) => res.json())
      .then((res) => setRows(res.shop));
  }, []);

  return (
    <div style={{ width: "90%", margin: "auto", paddingTop: "50px" }}>
      <DataGrid
        rows={rows}
        columns={columns}
        initialState={{
          pagenation: {
            paginationMode: { page: 0, pageSize: 5 },
          },
        }}
        pageSizeOptions={[5]}
        slots={{
          toolbar: EditToolbar,
        }}
        slotProps={{
          toolbar: { setRows },
        }}
        disableRowSelectionOnClick={true}
        editMode="row"
      />
    </div>
  );
};
