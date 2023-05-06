import WorldIcon from '../components/world-icon'

function MapLoadingHolder({ className = "" }: { className?: string }) {
  return (
    <div className={className}>
      <WorldIcon className="icon" />
      <h1>Initializing the map</h1>
    </div>
  )
}

export default MapLoadingHolder
