import loadingCoffee from '@/assets/loading-coffee.svg'

type LoadingProps = {
  loading?: boolean
}

export default function Loading({ loading = true }: LoadingProps) {
  if (!loading) return null
  return (
    <div className="fixed inset-0 z-[1000] grid place-items-center bg-white">
      <img src={loadingCoffee} alt="loading" className="w-28 animate-pulse" />
    </div>
  )
}
