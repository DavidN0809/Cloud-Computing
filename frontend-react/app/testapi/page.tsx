'use client';



import React, { useState } from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import TextField from '@mui/material/TextField';
import DialogActions from '@mui/material/DialogActions';




export default function TaskAction() {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [assignedTo, setAssignedTo] = useState('');
  const [status, setStatus] = useState('');
  const [hours, setHours] = useState('');
  const [start_date, setStartDate] = useState('');
  const [end_date, setEndDate] = useState('');


  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleCreate = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    const taskData = {
      title,
      description,
      assignedTo,
      status,
      hours,
      start_date,
      end_date
    };

    try {
      const response = await fetch('http://localhost:8002/tasks/create', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(taskData)
      });
      console.log(JSON.stringify(taskData));
      const data = await response.json(); // Moved outside of condition to simplify logic
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      
      console.log('Task created', data);

    } catch (error) {
      console.error('There was a problem with the fetch operation:', error);

    }
  };

  return (
    <React.Fragment>
      <div>
        <Button variant="outlined" color="primary" onClick={handleClickOpen}>
          Create Task
        </Button>
        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>Create New Task</DialogTitle>
          <DialogContent>
            <TextField autoFocus margin="dense" id="title" label="Task Name" type="text" fullWidth variant="standard" value={title} onChange={(e) => setTitle(e.target.value)} />
            <TextField margin="dense" id="description" label="Task Description" type="text" fullWidth multiline variant="standard" value={description} onChange={(e) => setDescription(e.target.value)} />
            <TextField margin="dense" id="assigned_to" label="Assigned To" type="text" fullWidth variant="standard" value={assignedTo} onChange={(e) => setAssignedTo(e.target.value)} />
            <TextField margin="dense" id="status" label="Status" type="text" fullWidth variant="standard" value={status} onChange={(e) => setStatus(e.target.value)} />
            <TextField margin="dense" id="hours" label="Hours" type="number" fullWidth variant="standard" value={hours} onChange={(e) => setHours(e.target.value)} />
            <TextField margin="dense" id="start_date" label="start date" type="number" fullWidth variant="standard" value={start_date} onChange={(e) => setStartDate(e.target.value)} />
            <TextField margin="dense" id="end_date" label="end date" type="number" fullWidth variant="standard" value={end_date} onChange={(e) => setEndDate(e.target.value)} />
          </DialogContent>
          <DialogActions>
            <Button onClick={handleClose} color="secondary">
              Cancel
            </Button>
            <Button onClick={handleCreate} color="primary">
              Create
            </Button>
          </DialogActions>
        </Dialog>
      </div>
    </React.Fragment>
  );
}
