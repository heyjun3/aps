import { useEffect, useState } from "react";
import { DataGrid, GridActionsCellItem } from "@mui/x-data-grid";
import DeleteIcon from "@mui/icons-material/DeleteOutlined";
import config from "../config";

export const Shops = () => {
  const [rows, setRows] = useState([]);

  const handleDeleteClick = (id) => () => {
    setRows(rows.filter((row) => row.id !== id));
  };

  const columns = [
    { field: "id", headerName: "ID", width: 130 },
    { field: "site_name", headerName: "Site name", width: 130 },
    { field: "name", headerName: "Shop name", width: 130 },
    {
      field: "interval",
      headerName: "Interval",
      width: 90,
    },
    {
      field: "url",
      headerName: "URL",
      width: 300,
      aggregable: false,
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
        pageSizeOptions={[5, 10]}
      />
    </div>
  );
};
