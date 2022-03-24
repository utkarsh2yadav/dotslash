import {FormControl, Select, MenuItem, InputLabel} from '@mui/material';
export default function LangSelect() {
  return <FormControl size='medium'>
    <Select
      id="language"
      label="Language"
    
    >
      <MenuItem value={'Java'}>Java</MenuItem>
      <MenuItem value={'Javascript'}>Javascript</MenuItem>
      <MenuItem value={'Python'}>Python</MenuItem>
      <MenuItem value={'Golang'}>Golang</MenuItem>
    </Select>
  </FormControl>
}
