import React, { useState } from "react";
import AceEditor from "react-ace";

import "ace-builds/src-noconflict/mode-golang";
import "ace-builds/src-noconflict/theme-monokai";
import "ace-builds/src-noconflict/ext-language_tools";
import { Fab } from "@mui/material";
import { PlayArrow } from "@mui/icons-material";


export default function Editor(props) {

  let [code, setCode] = useState("")

  return <>
    <AceEditor
      value={code}
      enableBasicAutocompletion={true}
      enableLiveAutocompletion={true}
      width="100%"
      height="100%"
      mode="golang"
      theme="monokai"
      name="editor"
      editorProps={{ $blockScrolling: true }}
      onChange={(value) => {
        setCode(value)
      }}
    />
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

  </>
}
