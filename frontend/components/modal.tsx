import {
    Fragment,
    FunctionComponent,
    ReactNode,
} from 'react'

import {
    Dialog,
    Transition,
} from '@headlessui/react'

export const Modal: FunctionComponent<{
    isOpen: boolean
    children: ReactNode
}> = ({ isOpen, children }) => {
    return (
        <Transition appear show={isOpen} as={Fragment}>
            <Dialog as='div' className='relative z-10' onClose={() => { }}>
                <Transition.Child
                    as={Fragment}
                    enter='ease-out duration-300'
                    enterFrom='opacity-0'
                    enterTo='opacity-100'
                    leave='ease-in duration-200'
                    leaveFrom='opacity-100'
                    leaveTo='opacity-0'
                >
                    <div className="fixed inset-0 bg-black bg-opacity-80 backdrop-blur-sm" />
                </Transition.Child>
                <div className="fixed inset-0 overflow-y-auto">
                    {children}
                </div>
            </Dialog>
        </Transition>
    )
}