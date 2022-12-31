import React, {useState, useEffect} from 'react'
import BoxList from './BoxList'
import config from '../config'


const ApiFetch = () => {

  const [filenames, setFilenames] = useState([])
  useEffect(() => {
      fetch(`${config.fqdn}/api/list`, {method: 'GET', mode: "cors"})
      .then(res => res.json())
      .then(data => {
          setFilenames(data['list'])
      })
  }, [])

  const deleteFile = (value) => {
    setFilenames(
      filenames.filter((file) => (file !== value))
    )
    fetch(`${config.fqdn}/api/deleteFile/${value}`, {method: 'DELETE', mode: 'cors'})
    .then(res => res.json())
    .then(data => console.log(data))
  }

  return (
    <div>
        <BoxList filenames={filenames} deleteFunc={deleteFile}/>
    </div>
  )
}

export default ApiFetch