import { useEffect, useState } from "react";
import { DataGrid } from "@mui/x-data-grid";
import config from "../config";
import { Link } from "@mui/material";

const RenderSKU = (props) => {
  const sku = props.value;
  const url = `https://sellercentral-japan.amazon.com/inventory/ref=xx_invmgr_dnav_xx?search:${sku}`;
  return (
    <Link tabIndex={props.tabIndex} href={url} target="_blank">
      {sku}
    </Link>
  );
};

const RenderName = (props) => {
  const { name, url } = props.value;
  return (
    <Link tabIndex={props.tabIndex} href={url} target="_blank">
      {name}
    </Link>
  );
};

const RenderLowest = (props) => {
  const { price, point, percent } = props.value;
  return (
    <div>
      ¥ {price}
      <br /> {`${point}pts (${percent}%)`}
    </div>
  );
};

const RenderPoint = (props) => {
  const percentPoint = props.value;
  return <div>{percentPoint}%</div>;
};

const columns = [
  { field: "sellerSku", headerName: "SKU", width: 200, renderCell: RenderSKU },
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
    align: "center",
    headerAlign: "center",
  },
  {
    field: "percentPoint",
    headerName: "Point(%)",
    width: 90,
    editable: true,
    align: "center",
    headerAlign: "center",
    renderCell: RenderPoint,
  },
  {
    field: "lowest",
    headerName: "Lowest",
    width: 120,
    align: "center",
    headerAlign: "center",
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
    const fetchInventories = async () => {
      const res = await fetch(`${config.fqdn}/api/inventories`, {
        method: "GET",
        mode: "cors",
      });
      const body = await res.json();
      for (const [i, value] of Object.entries(body)) {
        value.id = i;
        value.itemName = {
          name: value.productName,
          url: `https://www.amazon.co.jp/dp/${value.asin}`,
        };
        value.price = value.CurrentPrice?.Amount;
        value.point = value.CurrentPrice?.Point;
        value.percentPoint = value.CurrentPrice?.PercentPoint;
        value.lowest = {
          price: value.LowestPrice?.Amount,
          point: value.LowestPrice?.Point,
          percent: value.LowestPrice?.PercentPoint,
        };
      }
      setRows(body);
    };
    fetchInventories();
  }, []);

  const processRowUpdate = async (newRow) => {
    const updateRow = { ...newRow, isNew: false };
    setRows(
      rows.map((row) =>
        row.id === "" ? updateRow : row.id === newRow.id ? updateRow : row
      )
    );
    console.warn(updateRow.sellerSku);
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
