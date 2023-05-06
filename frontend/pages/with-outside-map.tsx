import 'mapbox-gl/dist/mapbox-gl.css'

import * as React from 'react'

import MapboxMap from '../components/mapbox-map'
import mapbox from '../lib/map-wrapper'

const WithOutsideMap: React.FC = () => {
  const [viewport, setViewport] = React.useState({
    center: ["105.804817", "21.028511"],
    zoom: "9.00",
  })

  const {
    center: [lng, lat],
    zoom,
  } = viewport

  const onMapCreated = React.useCallback((map: mapboxgl.Map) => {
    mapbox.map = map

    mapbox.map.on("move", () => {
      setViewport({
        center: [
          mapbox.map.getCenter().lng.toFixed(4),
          mapbox.map.getCenter().lat.toFixed(4),
        ],
        zoom: mapbox.map.getZoom().toFixed(2),
      })
    })
  }, [])

  return (
    <div className="app-container">
      <div className="map-wrapper">
        <div className="viewport-panel">
          Longitude: {lng} | Latitude: {lat} | Zoom: {zoom}
        </div>
        <MapboxMap
          onCreated={onMapCreated}
          initialOptions={{ center: [+lng, +lat], zoom: +zoom }}
        />
      </div>
    </div>
  )
}

export default WithOutsideMap
