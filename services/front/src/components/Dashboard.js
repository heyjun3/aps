import React from "react"

import config from "../config"
// import OutlinedCard from "./OutlinedCard"

const Dashboard = () => {
    const [keepa, setKeepa] = React.useState({})
    const [mws, setMws] = React.useState({})

    React.useEffect(() => {
        fetch(`${config.fqdn}/api/counts`, {method: "GET", mode: "cors"})
        .then(res => res.json())
        .then(data => {
            setKeepa(data.keepa)
            setMws(data.mws)
        })
    }, [])
    return (
        <div>
            {keepa ? <h3>Keepa Table modified Count: {keepa.modified} / {keepa.total}</h3> : ""}
            {mws ? <h3>MWS Table price Count: {mws.price} / {mws.total}</h3> : ""}
            {mws ? <h3>MWS Table fee Count: {mws.fee} / {mws.total}</h3> : ""}
        </div>
    )
}

export default Dashboard