import * as React from 'react'

import Head from 'next/head'

import { Dialog } from '@headlessui/react'

import MapLoadingHolder from '../components/map-loading-holder'
import MapboxMap from '../components/mapbox-map'
import { Modal } from '../components/modal'

function App() {
  const [loading, setLoading] = React.useState(true)
  const handleMapLoading = () => setLoading(false)
  const [isOpen, setOpen] = React.useState(true)

  return (
    <>
      <Head>
        <title>Map</title>
      </Head>
      <div className="app-container">
        <div className="map-wrapper">
          <MapboxMap
            initialOptions={{ center: [38.0983, 55.7038] }}
            onLoaded={handleMapLoading}
          />
        </div>
        {loading && <MapLoadingHolder className="loading-holder" />}
      </div>
      <Modal isOpen={isOpen} >
        <div className='flex h-full items-center justify-center'>
          <Dialog.Panel className='w-full max-w-md transform overflow-hidden rounded-2xl bg-white p-10 text-left align-middle shadow-xl transition-all flex flex-col justify-center items-center'>
            <div className='w-[500px] h-[430px] flex flex-col justify-start items-start px-12'>
              <div className='w-full flex justify-center items-center text-3xl text-black font-black mb-3'>{'EBN TEAM'}</div>
              <label className='text-lg font-bold text-black'>{'Username (*)'}</label>
              <input placeholder='username' className='px-4 py-3 rounded-lg border border-solid border-black w-full mt-3' />
              <label className='mt-5 text-lg font-bold text-black'>{'Radius'}</label>
              <input placeholder='5 .km' className='px-4 py-3 rounded-lg border border-solid border-black w-full mt-3' />
              <label className='mt-5 text-lg font-bold text-black'>{'Your location (optional)'}</label>
              <div className='w-full flex flex-row mt-3 gap-2'>
                <input placeholder='Longitude' className='w-full px-4 py-3 rounded-lg border border-solid border-black ' />
                <input placeholder='Latitude' className='w-full px-4 py-3 rounded-lg border border-solid border-black ' />
              </div>
              <button onClick={() => setOpen(false)} className='py-3 w-full mt-8 rounded-lg bg-black text-white font-bold text-lg hover:bg-gray-600'>{'Join Application'}</button>
            </div>
          </Dialog.Panel>
        </div>
      </Modal>
    </>
  )
}

export default App
