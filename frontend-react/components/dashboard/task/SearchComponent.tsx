import React, { useState } from 'react';
import TextField from '@mui/material/TextField';
import Switch from '@mui/material/Switch';
import FormControlLabel from '@mui/material/FormControlLabel';
import Button from '@mui/material/Button';

function SearchComponent() {
  const [searchType, setSearchType] = useState('taskid'); // 默认按 taskid 搜索
  const [searchText, setSearchText] = useState('');

  const handleSwitchChange = (event:React.ChangeEvent<HTMLInputElement>) => {
    setSearchType(event.target.checked ? 'userid' : 'taskid');
  };

  const handleSearch = () => {
    // 实现搜索逻辑
    console.log(`Searching ${searchType} for: ${searchText}`);
  };

  return (
    <div>
      <FormControlLabel
        control={
          <Switch
            checked={searchType === 'userid'}
            onChange={handleSwitchChange}
            name="searchTypeSwitch"
          />
        }
        label={searchType === 'taskid' ? 'Search by Task ID' : 'Search by User ID'}
      />
      <TextField
        label={searchType === 'taskid' ? 'Task ID' : 'User ID'}
        variant="outlined"
        value={searchText}
        onChange={(e) => setSearchText(e.target.value)}
      />
      <Button onClick={handleSearch} variant="contained" color="primary">
        Search
      </Button>
    </div>
  );
}

export default SearchComponent;
