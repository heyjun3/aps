import { useEffect, useState } from "react";
import { DataGrid, GridToolbarContainer } from "@mui/x-data-grid";
import Button from "@mui/material/Button";
import { Link } from "@mui/material";
import SaveIcon from "@mui/icons-material/Save";

import config from "../config";

const Toolbar = (props) => {
  const {
    rows,
    setRows,
    updateRows,
    setUpdateRows,
    tmpRows,
    setTmpRows,
    isShowUpdateRows,
    setIsShowUpdateRows,
  } = props;

  const togleRows = () => {
    console.warn(updateRows)
    if (isShowUpdateRows) {
      setRows([...tmpRows]);
      setIsShowUpdateRows(!isShowUpdateRows);
      return
    }
    setRows([...updateRows]);
    setTmpRows([...rows]);
    setIsShowUpdateRows(!isShowUpdateRows);
  };
  const refreshInventories = async () => {
    setRows([]);
    const sleep = (msec) => new Promise(resolve => setTimeout(resolve, msec))
    await sleep(10000)
    await fetchInventories(setRows)
  };
  const cancel = () => {
    setUpdateRows([])
    setIsShowUpdateRows(false)
    fetchInventories(setRows)
  }
  const saveAll = () => {
    console.log(updateRows);
  };

  return (
    <GridToolbarContainer>
      <Button
        color="primary"
        startIcon={<SaveIcon />}
        onClick={refreshInventories}
      >
        Refresh inventories
      </Button>
      <Button color="primary" startIcon={<SaveIcon />} onClick={togleRows}>
        {isShowUpdateRows ? "show all" : "show updated"}
      </Button>
      <Button color="primary" startIcon={<SaveIcon />} onClick={cancel}>
        Cancel
      </Button>
      <Button color="primary" startIcon={<SaveIcon />} onClick={saveAll}>
        Save all
      </Button>
    </GridToolbarContainer>
  );
};

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
];

const fetchInventories = async (set) => {
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
  set(body);
};

export const Items = () => {
  const [rows, setRows] = useState([]);
  const [updateRows, setUpdateRows] = useState([]);
  const [tmpRows, setTmpRows] = useState([]);
  const [isShowUpdateRows, setIsShowUpdateRows] = useState(false);

  useEffect(() => {
    fetchInventories(setRows);
  }, []);

  const processRowUpdate = async (newRow, oldRow) => {
    const isSamePriceAndPercentPoint = (row, old) => {
      return (
        Number(row.price) === Number(old.price) &&
        Number(row.percentPoint) === Number(old.percentPoint)
      );
    };
    const updateRow = { ...newRow, isNew: false };
    setRows(
      rows.map((row) =>
        row.id === "" ? updateRow : row.id === newRow.id ? updateRow : row
      )
    );
    if (isSamePriceAndPercentPoint(newRow, oldRow)) {
      return updateRow;
    }
    setUpdateRows([...updateRows, updateRow]);
    // const res = await fetch(`${config.fqdn}/api/price/${updateRow.sellerSku}`, {
    //   method: "POST",
    //   mode: "cors",
    //   headers: {
    //     "Content-Type": "application/json",
    //   },
    //   body: JSON.stringify({
    //     price: Number(updateRow.price),
    //     percentPoint: Number(updateRow.percentPoint),
    //   }),
    // });
    // const body = await res.json();
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
        slots={{
          toolbar: Toolbar,
        }}
        slotProps={{
          toolbar: {
            rows,
            setRows,
            updateRows,
            setUpdateRows,
            tmpRows,
            setTmpRows,
            isShowUpdateRows,
            setIsShowUpdateRows,
          },
        }}
      />
    </div>
  );
};
