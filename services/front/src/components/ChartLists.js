import React, { useEffect, useState } from "react";
import { useLocation, Navigate } from "react-router-dom";
import Pagination from "@mui/material/Pagination";
import Stack from "@mui/material/Stack";
import RenderLineChart from "./LineChart";
import config from "../config";
import { PageBox } from "./PageBox";

const filenameNumber = 2;

const ChartLists = () => {
  const [products, setProducts] = useState([]);
  const [redirect, setRedirect] = useState(false);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(100);
  const [maxPage, setMaxPage] = useState(0);

  let location = useLocation();
  const filename = location.pathname.split("/")[filenameNumber];

  useEffect(() => {
    const params = { page, limit };
    const query = new URLSearchParams(params);
    fetch(`${config.fqdn}/api/chart_list/${filename}?${query}`, {
      method: "GET",
      mode: "cors",
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.status === "error") {
          console.log("error");
          setRedirect(true);
        }
        setProducts(data.chart_data);
        setMaxPage(data.max_page);
      });
  }, [page, filename, limit]);

  const handleChange = (event, value) => {
    setPage(value);
    window.scrollTo({
      top: 0,
      behavior: "auto",
    });
  };

  if (redirect) {
    return <Navigate to="/list" />;
  }

  return (
    <div className="chartLists">
      <PageBox setLimit={setLimit} />
      {products.map((product) => {
        return (
          <RenderLineChart
            key={product.asin}
            data={product.Chart.data}
            title={product.title}
            jan={product.jan}
            asin={product.asin}
            url={product.url}
            diffCount={product.Chart.diff_count}
            diffCountMA7={product.Chart.diff_count_ma7}
          />
        );
      })}
      <Stack className="stack" spacing={2}>
        <Pagination
          className="pagination"
          count={maxPage}
          page={page}
          onChange={handleChange}
        />
      </Stack>
    </div>
  );
};

export default ChartLists;
