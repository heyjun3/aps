import * as React from 'react';
import Paper from '@mui/material/Paper';
import InputBase from '@mui/material/InputBase';
import Divider from '@mui/material/Divider';
import IconButton from '@mui/material/IconButton';
import SearchIcon from '@mui/icons-material/Search';
import RenderLineChart from './LineChart';



export default function CustomizedInputBase() {
  const [data, setData] = React.useState([])
  const [asin, setAsin] = React.useState('')

  const handleClick = () => {
      const value = asin
      if (value){
        console.log(value)
        fetch(`http://${window.location.hostname}:5000/search/${value}`, {method: 'GET', mode: 'cors'})
        .then(res => res.json())
        .then(data => setData(data))
      } else {
        setData([])
      }
  }

  const handleInputChange = () => evt => {
    const asin = evt.target.value
    setAsin(asin)
    }

  return (
    <div>
    <Paper
      component="form"
      sx={{ p: '2px 4px', display: 'flex', alignItems: 'center', width: '50%' , margin: 'auto', mt: '50px'}}
    >
      <InputBase
        sx={{ ml: 1, flex: 1 }}
        placeholder="Search by Amazon, Enter ASIN Code"
        inputProps={{ 'aria-label': 'search amazon' }}
        value={asin}
        onChange={handleInputChange()}
      />
      <Divider sx={{ height: 28, m: 0.5 }} orientation="vertical" />
      <IconButton type="button" sx={{ p: '10px' }} aria-label="search" onClick={handleClick}>
        <SearchIcon />
      </IconButton>
    </Paper>
      {data.length !== 0 && <RenderLineChart data={data.chart_data} title={data.title} jan={data.jan} asin={data.asin}/>}
    </div>
    
  );
}
