import React from 'react'
import {LineChart, Line, CartesianGrid, XAxis, YAxis, ResponsiveContainer, Legend} from 'recharts';


const RenderLineChart = (props) => {
  return (
    <div className="renderLineChart">
      <h3 className="title">{props.title}</h3>
      <ResponsiveContainer className="LineChart" width='70%' height={400}>
      <LineChart data={props.data} >
        <YAxis yAxisId={1} type="number" domain={["dataMin - 1000", "dataMax + 1000"]}/>
        <YAxis yAxisId={2} type="number" domain={["dataMin - 1000", "dataMax + 2000"]} orientation="right"/>
        <Line yAxisId={1} strokeWidth={3} type="monotone" dataKey="price" stroke="#8884d8" dot={false} isAnimationActive={false}/>
        <Line yAxisId={2} strokeWidth={3} type="monotone" dataKey="rank" stroke="#82ca9d" dot={false} isAnimationActive={false}/>
        <CartesianGrid stroke="#ccc" strokeDasharray="5 5" />
        <XAxis dataKey="date" />
        <Legend />
      </LineChart>
      </ResponsiveContainer>
      <div className="asinBox">
        <h3 className="asin">JAN : <a href={props.url} target="_blank" rel="noreferrer">{props.jan}</a> ASIN : <a href={`https://www.amazon.co.jp/dp/${props.asin}`} target="_blank" rel="noreferrer">{props.asin}</a></h3></div>
    </div>
  )
}

export default RenderLineChart