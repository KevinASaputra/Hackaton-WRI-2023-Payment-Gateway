interface MainContainerProps {
  children?: React.ReactNode
}

export const MainContainer = ({ children }: MainContainerProps) => {
  return (
    <div className="flex mx-auto max-w-7xl min-h-screen bg-slate-700 justify-center">
      {children}
    </div>
  )
}