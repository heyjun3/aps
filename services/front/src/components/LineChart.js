import React from "react";
import {
  LineChart,
  Line,
  CartesianGrid,
  XAxis,
  YAxis,
  ResponsiveContainer,
  Legend,
} from "recharts";
import { Box, Typography as P } from "@mui/material";

const RenderLineChart = (props) => {
  return (
    <div className="renderLineChart">
      <h3 className="title">{props.title}</h3>
      <ResponsiveContainer className="LineChart" width="70%" height={400}>
        <LineChart data={props.data}>
          <YAxis
            yAxisId={1}
            type="number"
            domain={["dataMin - 1000", "dataMax + 1000"]}
          />
          <YAxis
            yAxisId={2}
            type="number"
            domain={["dataMin - 1000", "dataMax + 2000"]}
            orientation="right"
          />
          <Line
            yAxisId={1}
            strokeWidth={3}
            type="monotone"
            dataKey="price"
            stroke="#8884d8"
            dot={false}
            isAnimationActive={false}
          />
          <Line
            yAxisId={2}
            strokeWidth={3}
            type="monotone"
            dataKey="rank"
            stroke="#82ca9d"
            dot={false}
            isAnimationActive={false}
          />
          <Line
            yAxisId={2}
            strokeWidth={3}
            type="monotone"
            dataKey="rank_ma7"
            stroke="red"
            dot={false}
            isAnimationActive={false}
          />
          <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
          <XAxis dataKey="date" />
          <Legend />
        </LineChart>
      </ResponsiveContainer>
      <Box
        width={"70%"}
        display={"flex"}
        margin={"0 auto"}
        gap={"10px"}
        justifyContent={"end"}
      >
        <P fontWeight={"bold"}>DropsMA7: {props.diffCountMA7}</P>
        <P className="asin" fontWeight={"bold"}>
          JAN :
          <a href={props.url} target="_blank" rel="noreferrer">
            {props.jan}
          </a>
        </P>
        <P className="asin" fontWeight={"bold"}>
          ASIN :
          <a
            href={`https://www.amazon.co.jp/dp/${props.asin}`}
            target="_blank"
            rel="noreferrer"
          >
            {props.asin}
          </a>
        </P>
      </Box>
    </div>
  );
};

export default RenderLineChart;
