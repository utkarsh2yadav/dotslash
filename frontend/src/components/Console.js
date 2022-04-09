import { FitAddon } from 'xterm-addon-fit'
import { useEffect } from 'react'
import ResizeObserver from 'react-resize-observer'
import "xterm/css/xterm.css"

export default function Console(props) {
  let fitAddon = new FitAddon()
  let xterm = props.xterm

  useEffect(() => {
    xterm.open(document.getElementById("terminal"))
    xterm.loadAddon(fitAddon)
    fitAddon.fit()
  })

  return (
    <>
      <div id='terminal' style={{ width: "100%", height: "100%" }} />
      <ResizeObserver onResize={() => fitAddon.fit()} />
    </>
  )
}

