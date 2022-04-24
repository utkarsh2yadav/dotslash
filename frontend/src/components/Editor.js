import React, { useState } from "react"
import AceEditor from "react-ace"
import { Fab } from "@mui/material"
import { PlayArrow, StopRounded } from "@mui/icons-material"

import "ace-builds/src-noconflict/mode-c_cpp"
import "ace-builds/src-noconflict/mode-golang"
import "ace-builds/src-noconflict/mode-java"
import "ace-builds/src-noconflict/mode-javascript"
import "ace-builds/src-noconflict/mode-python"
import "ace-builds/src-noconflict/mode-typescript"
import "ace-builds/src-noconflict/theme-one_dark"

import "ace-builds/src-noconflict/ext-language_tools"
import beautify from "ace-builds/src-noconflict/ext-beautify"

import "ace-builds/src-noconflict/snippets/c_cpp"
import "ace-builds/src-noconflict/snippets/golang"
import "ace-builds/src-noconflict/snippets/java"
import "ace-builds/src-noconflict/snippets/javascript"
import "ace-builds/src-noconflict/snippets/python"
import "ace-builds/src-noconflict/snippets/typescript"


export default function Editor(props) {

  let [code, setCode] = useState("")
  let ws
  let input = ""

  props.xterm.onKey(e => {
    if (ws.readyState === ws.OPEN) {
      if (e.domEvent.keyCode === 13) {
        ws.send(JSON.stringify({ input: input }))
        props.xterm.write("\r\n")
        input = ""
      } else if (e.domEvent.keyCode === 8) {
        if (input.length > 0) {
          input = input.slice(0, -1)
        }
        props.xterm.write("\b \b")
      } else {
        input += e.key
        props.xterm.write(e.key)
      }
    }
  })


  return <>
    <AceEditor
      value={code}
      enableBasicAutocompletion={true}
      enableLiveAutocompletion={true}
      width="100%"
      height="100%"
      mode="golang"
      theme="one_dark"
      name="editor"
      editorProps={{ $blockScrolling: true }}
      onChange={(value) => {
        setCode(value)
      }}
      commands={beautify.commands}
    />
    Ctrl + Shift + B to beautify
    <Fab color='secondary' variant="extended" onClick={(_) => {
      let ws = new WebSocket("ws://localhost:8080/ws/golang")
      ws.onopen = () => {
        ws.send(JSON.stringify({ "code": code }))
      }
      ws.onmessage = (event) => {
        props.xterm.write(event.data)
      }
    }}>
      <PlayArrow sx={{ mr: 1 }} />
      Run
    </Fab>

    <Fab color='warning' variant="extended" onClick={(_) => {
      if (ws && ws.readyState === ws.OPEN) {
        ws.send(JSON.stringify({ interrupt: true }))
      }
    }}>
      <StopRounded sx={{ mr: 1 }} />
      Stop
    </Fab>
  </>
}
