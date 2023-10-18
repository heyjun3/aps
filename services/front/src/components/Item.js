import { useEffect, useState } from "react";
import {
  DataGrid,
  GridActionsCellItem,
  GridToolbarContainer,
} from "@mui/x-data-grid";
import DeleteIcon from "@mui/icons-material/DeleteOutlined";
import Button from "@mui/material/Button";
import AddIcon from "@mui/icons-material/Add";
import config from "../config";
import { Link } from "@mui/material";

const EditToolbar = (props) => {
  const { setRows } = props;

  const handleClick = () => {
    setRows((oldRows) => [
      {
        id: "",
        site_name: "",
        name: "",
        interval: "",
        url: "",
        isNew: true,
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

const RenderName = (props) => {
  const { name, url } = props.value;
  return (
    <Link tabIndex={props.tabIndex} href={url} target="_blank">
      {name}
    </Link>
  );
};

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
      },
    ];
    setRows(rows);
  }, []);

  // useEffect(() => {
  //   fetch(`${config.fqdn}/api/shops`, { method: "GET", mode: "cors" })
  //     .then((res) => res.json())
  //     .then((res) => setRows(res.shop));
  // }, []);

  const handleDeleteClick = (id) => async () => {
    setRows(rows.filter((row) => row.id !== id));
    const reqBody = { ids: [id] };
    const res = await fetch(`${config.fqdn}/api/shops`, {
      method: "DELETE",
      mode: "cors",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(reqBody),
    });
    const data = await res.json();
    if (data != null) {
      console.warn(data);
    }
  };

  const processRowUpdate = async (newRow) => {
    const updateRow = { ...newRow, isNew: false };
    setRows(
      rows.map((row) =>
        row.id === "" ? updateRow : row.id === newRow.id ? updateRow : row
      )
    );

    const reqBody = { shop: [updateRow] };
    const res = await fetch(`${config.fqdn}/api/shops`, {
      method: "POST",
      mode: "cors",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(reqBody),
    });
    const body = await res.json();
    if (body != null) {
      console.warn(body);
      return;
    }
    return updateRow;
  };

  const columns = [
    { field: "sku", headerName: "SKU", width: 200 },
    {
      field: "itemName",
      headerName: "Name",
      width: 100,
      flex: 1,
      renderCell: RenderName,
    },
    { field: "price", headerName: "Price", width: 100, editable: true },
    { field: "point", headerName: "Point", width: 100, editable: true },
    // {
    //   field: "actions",
    //   type: "actions",
    //   headerName: "Actions",
    //   width: 100,
    //   cellClassName: "actions",
    //   getActions: ({ id }) => {
    //     return [
    //       <GridActionsCellItem
    //         icon={<DeleteIcon />}
    //         label="Delete"
    //         onClick={handleDeleteClick(id)}
    //         color="inherit"
    //       />,
    //     ];
    //   },
    // },
  ];

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
        slots={{
          toolbar: EditToolbar,
        }}
        slotProps={{
          toolbar: { setRows },
        }}
        disableRowSelectionOnClick={true}
        editMode="row"
        processRowUpdate={processRowUpdate}
      />
    </div>
  );
};
