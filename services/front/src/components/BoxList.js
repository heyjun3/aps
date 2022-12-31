import * as React from 'react';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemText from '@mui/material/ListItemText';
import IconButton from '@mui/material/IconButton';
import DeleteIcon from '@mui/icons-material/Delete';
import {Link} from 'react-router-dom';


export default function BoxList(props) {

  return (
    <List sx={{ width: '100%', maxWidth: 500, bgcolor: 'background.paper', margin: 'auto', mt: '10px' }}>
      {props.filenames.map((value) => {
        const labelId = `checkbox-list-label-${value}`;

        return (
          <ListItem
            key={value}
            secondaryAction={
              <IconButton edge="end" aria-label="comments" onClick={() => props.deleteFunc(value)}>
                <DeleteIcon />
              </IconButton>
            }
            disablePadding
          >
            <ListItemButton component={Link} to={`/chartList/${value}`} role={undefined} dense target="_blank" >
              <ListItemText align="center" id={labelId} primary={value} />
            </ListItemButton>
          </ListItem>
        );
      })}
    </List>
  );
}
