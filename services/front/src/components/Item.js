import { useEffect, useState } from "react";
import { DataGrid } from "@mui/x-data-grid";
import config from "../config";
import { Link } from "@mui/material";

const RenderName = (props) => {
  const { name, url } = props.value;
  return (
    <Link tabIndex={props.tabIndex} href={url} target="_blank">
      {name}
    </Link>
  );
};

const RenderLowest = (props) => {
  const { price, point } = props.value;
  return (
    <div>
      price: {price}
      <br />
      point: {point}
    </div>
  );
};

const columns = [
  { field: "sku", headerName: "SKU", width: 200 },
  {
    field: "itemName",
    headerName: "Name",
    width: 100,
    flex: 1,
    renderCell: RenderName,
    aggregable: false,
  },
  {
    field: "price",
    headerName: "Price",
    width: 90,
    editable: true,
  },
  { field: "point", headerName: "Point", width: 90, editable: true },
  {
    field: "lowest",
    headerName: "Lowest",
    width: 120,
    renderCell: RenderLowest,
  },
  {
    field: "update",
    headerName: "Update",
    width: 120,
  },
];

export const Items = () => {
  const [rows, setRows] = useState([]);

  useEffect(() => {
    const rows = [
      {
        id: "id",
        sku: "4710483940767-N-28480-20231016",
        itemName: {
          name: "ASRock マザーボード Z790 Pro RS Intel 第12世代 ・ 13世代 CPU ( LGA1700 )対応 Z790チップセット DDR5 ATX マザーボード 【国内正規代理店品】",
          url: "https://www.amazon.co.jp/dp/B0BJDZ4YZR?ref=myi_title_dp",
        },
        price: 30000,
        point: 3000,
        lowest: { price: 300000, point: 3000 },
        update: "2023/01/01",
      },
    ];
    setRows(rows);
  }, []);

  const processRowUpdate = async (newRow) => {
    const updateRow = { ...newRow, isNew: false };
    setRows(
      rows.map((row) =>
        row.id === "" ? updateRow : row.id === newRow.id ? updateRow : row
      )
    );

    // const reqBody = { shop: [updateRow] };
    // const res = await fetch(`${config.fqdn}/api/shops`, {
    //   method: "POST",
    //   mode: "cors",
    //   headers: { "Content-Type": "application/json" },
    //   body: JSON.stringify(reqBody),
    // });
    // const body = await res.json();
    // if (body != null) {
    //   console.warn(body);
    //   return;
    // }
    return updateRow;
  };

  return (
    <div style={{ width: "95%", margin: "auto", paddingTop: "50px" }}>
      <DataGrid
        rows={rows}
        columns={columns}
        initialState={{
          pagenation: {
            paginationMode: { page: 0, pageSize: 5 },
          },
        }}
        pageSizeOptions={[5]}
        disableRowSelectionOnClick={true}
        editMode="row"
        processRowUpdate={processRowUpdate}
      />
    </div>
  );
};
