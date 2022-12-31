import React, {useEffect, useState} from 'react'
import {useLocation, Navigate} from 'react-router-dom'
import Pagination from '@mui/material/Pagination'
import Stack from '@mui/material/Stack'
import RenderLineChart from './LineChart'
import config from '../config'

const filenameNumber = 2

const ChartLists = () => {
  const [products, setProducts] = useState([])
  const [redirect, setRedirect] = useState(false)
  const [page, setPage] = useState(1)
  const [maxPage, setMaxPage] = useState(0)

  let location = useLocation()
  const filename = location.pathname.split('/')[filenameNumber]

  useEffect(() => {
    fetch(`${config.fqdn}/api/chart_list/${filename}?page=${page}`, {method: 'GET', mode: 'cors'})
    .then(res => res.json()) 
    .then(data => {
      if (data.status === "error"){
        console.log('error')
        setRedirect(true)
      }
      setProducts(data.chart_data)
      setMaxPage(data.max_page)
    })
  }, [page, filename])

  const handleChange = (event, value) => {
    setPage(value)
    window.scrollTo({
      top: 0,
      behavior: "auto",
    })
  }
  
  if (redirect){
    return <Navigate to="/list"/>
  }

  return (
    <div className="chartLists">
      {products.map((product) => { 
        return <RenderLineChart key={product.asin} data={product.data} title={product.title} jan={product.jan} asin={product.asin} url={product.url}/>
      })}
      <Stack className="stack" spacing={2}>
        <Pagination className="pagination" count={maxPage} page={page} onChange={handleChange}/>
      </Stack>
    </div>
  )
}

export default ChartLists