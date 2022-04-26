import { FormControl, Select, MenuItem } from '@mui/material';


export default function LangSelect() {
  return (
    <FormControl sx={{ m: 1, minWidth: 80 }}>
      <Select
        id="language"
        label="Language"
        value={'Java'}
        autoWidth
      >
        <MenuItem value={'Java'}>Java</MenuItem>
        <MenuItem value={'Javascript'}>Javascript</MenuItem>
        <MenuItem value={'Python'}>Python</MenuItem>
        <MenuItem value={'Golang'}>Golng</MenuItem>
      </Select>
      
    </FormControl>
  )
}
