import { useEffect, useRef } from 'react'
import { AlertTriangle, Trash2, RotateCcw, XCircle } from 'lucide-react'

type ModalType = 'delete' | 'restore' | 'force-delete' | 'warning'

interface ConfirmModalProps {
  isOpen: boolean
  title: string
  message: string
  type?: ModalType
  confirmText?: string
  onConfirm: () => void
  onCancel: () => void
}

const modalConfig: Record<ModalType, { icon: React.ReactNode; bgColor: string; textColor: string; buttonColor: string }> = {
  delete: {
    icon: <Trash2 className="w-6 h-6" />,
    bgColor: 'bg-red-100',
    textColor: 'text-red-600',
    buttonColor: 'bg-red-600 hover:bg-red-700',
  },
  restore: {
    icon: <RotateCcw className="w-6 h-6" />,
    bgColor: 'bg-green-100',
    textColor: 'text-green-600',
    buttonColor: 'bg-green-600 hover:bg-green-700',
  },
  'force-delete': {
    icon: <XCircle className="w-6 h-6" />,
    bgColor: 'bg-red-100',
    textColor: 'text-red-600',
    buttonColor: 'bg-red-600 hover:bg-red-700',
  },
  warning: {
    icon: <AlertTriangle className="w-6 h-6" />,
    bgColor: 'bg-yellow-100',
    textColor: 'text-yellow-600',
    buttonColor: 'bg-yellow-600 hover:bg-yellow-700',
  },
}

export default function ConfirmModal({
  isOpen,
  title,
  message,
  type = 'delete',
  confirmText,
  onConfirm,
  onCancel,
}: ConfirmModalProps) {
  const modalRef = useRef<HTMLDivElement>(null)
  const config = modalConfig[type]

  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel()
    }
    if (isOpen) {
      document.addEventListener('keydown', handleEscape)
      document.body.style.overflow = 'hidden'
    }
    return () => {
      document.removeEventListener('keydown', handleEscape)
      document.body.style.overflow = 'unset'
    }
  }, [isOpen, onCancel])

  if (!isOpen) return null

  const getConfirmText = () => {
    if (confirmText) return confirmText
    switch (type) {
      case 'delete':
        return 'Delete'
      case 'restore':
        return 'Restore'
      case 'force-delete':
        return 'Delete Permanently'
      case 'warning':
        return 'Confirm'
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      <div
        className="fixed inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onCancel}
      />
      <div
        ref={modalRef}
        className="relative bg-white rounded-2xl shadow-2xl p-6 w-full max-w-sm mx-4 transform transition-all"
      >
        <div className="text-center">
          <div className={`mx-auto w-12 h-12 ${config.bgColor} rounded-full flex items-center justify-center mb-4 ${config.textColor}`}>
            {config.icon}
          </div>
          <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
          <p className="text-gray-500 text-sm mb-6">{message}</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={onCancel}
            className="flex-1 px-4 py-2.5 text-gray-700 bg-gray-100 rounded-xl font-medium hover:bg-gray-200 transition-colors"
          >
            Cancel
          </button>
          <button
            onClick={onConfirm}
            className={`flex-1 px-4 py-2.5 text-white rounded-xl font-medium transition-colors ${config.buttonColor}`}
          >
            {getConfirmText()}
          </button>
        </div>
      </div>
    </div>
  )
}
