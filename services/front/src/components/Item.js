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
    if (isShowUpdateRows) {
      setRows(
        tmpRows.map((row) => {
          const update = updateRows.find((up) => up.id === row.id);
          return update ? update : row;
        })
      );
      setIsShowUpdateRows(!isShowUpdateRows);
      return;
    }
    setRows(updateRows);
    setTmpRows(rows);
    setIsShowUpdateRows(!isShowUpdateRows);
  };

  const refreshInventories = async () => {
    setRows([]);
    await fetch(`${config.fqdn}/api/inventory/refresh`, {
      method: "POST",
      mode: "cors",
    });
    await fetch(`${config.fqdn}/api/price/refresh`, {
      method: "POST",
      mode: "cors",
    });
    await fetch(`${config.fqdn}/api/lowest-price/refresh`, {
      method: "POST",
      mode: "cors",
    });
    await fetchInventories(setRows);
  };
  const cancel = async () => {
    setUpdateRows([]);
    setIsShowUpdateRows(false);
    await fetchInventories(setRows);
  };
  const saveAll = async () => {
    if (updateRows.length === 0) {
      return
    }

    const body = updateRows.map((row) => ({
      sku: row.sellerSku,
      price: Number(row.price),
      percentPoint: Number(row.percentPoint),
    }));

    await fetch(`${config.fqdn}/api/price/update`, {
      method: "POST",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(body),
    });

    setRows([])
    setUpdateRows([])
    setTmpRows([])
    await fetchInventories(setRows);
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
    setUpdateRows([
      ...updateRows.filter((row) => row.id !== updateRow.id),
      updateRow,
    ]);
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
