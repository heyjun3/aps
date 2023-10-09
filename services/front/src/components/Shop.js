import { useEffect, useState } from "react";
import { DataGrid } from "@mui/x-data-grid";
import config from '../config'

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
  },
];

export const Shops = () => {
  const [rows, setRows] = useState([])
  useEffect(() => {
    fetch(`${config.fqdn}/api/shops`, {method: "GET", mode: "cors"})
    .then(res => res.json())
    .then(res => setRows(res.shop))
  }, [])
  return (
    <div style={{ width: "80%", margin: "auto", paddingTop: "50px" }}>
      <DataGrid
        rows={rows}
        columns={columns}
        initialState={{
          pagenation: {
            paginationMode: { page: 0, pageSize: 5 },
          },
        }}
        pageSizeOptions={[5, 10]}
        checkboxSelection
      />
    </div>
  );
};
