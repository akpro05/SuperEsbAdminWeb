import React, { useState, useEffect, useRef } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { TablePagination } from '@mui/material';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, InputAdornment } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import PictureAsPdfIcon from '@mui/icons-material/PictureAsPdf';
import SearchIcon from '@mui/icons-material/Search';
import { Link } from 'react-router-dom';
// Import the daterangepicker CSS
import dayjs from 'dayjs';
import 'dayjs/locale/en';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap-daterangepicker/daterangepicker.css';
import moment from 'moment';
import $ from 'jquery';
import 'popper.js';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import 'bootstrap-daterangepicker/daterangepicker.js';
import CopyToClipboard from 'react-copy-to-clipboard';
import * as XLSX from 'xlsx';
import jsPDF from 'jspdf';
import 'jspdf-autotable';

//for making HTTP requests
import axios from 'axios';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import SecondaryAction from '../../../ui-component/cards/CardSecondaryAction';
import { CREATE_USER_URL, SEARCH_USER_URL } from '../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import HomeIcon from '@mui/icons-material/Home';

// assets
import LinkIcon from '@mui/icons-material/Link';
import '../../../assets/scss/style.css'

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================|| TABLER ICONS ||============================= //

const TablerIcons = () => {
  const navigate = useNavigate();
  const [errors, setErrors] = useState({});
  const [rows, setRows] = useState([]);
  const { t, i18n } = useTranslation();
  const today = dayjs();
  const lastMonthStart = today.subtract(1, 'month').startOf('month');
  const lastMonthEnd = today.subtract(1, 'month').endOf('month');
  const [customDateRangeText, setCustomDateRangeText] = useState(
    `${lastMonthStart.format('MM/DD/YYYY')} - ${today.format('MM/DD/YYYY')}`
  );
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  // Data to copy or export to xlsx or pdf in table
  const [copied, setCopied] = useState(false);
  const [tableData, setTableData] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');

  // Add the exportToExcel function
  const exportToExcel = () => {
    // Customize the data before exporting
    const customizedTableData = tableData.map((row) => ({
      'Id': row.Id,
      'Timestamp': row.Timestamp,
      'Full Name': row.FullName,
      'Email': row.Email,
      'Status': row.Status,
      'Mobile Number': row.Mobile,
      'Address': row.Address,

    }));

    const ws = XLSX.utils.json_to_sheet(customizedTableData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    XLSX.writeFile(wb, 'SystemUserData.xlsx');
  };

  // Add the exportToPdf function
  const exportToPDF = () => {
    const doc = new jsPDF();
    doc.text('Table Data', 10, 10);

    // Check the table ID used in your HTML, it should match your table's ID
    doc.autoTable({ startY: 20, html: '#tableId' });

    doc.save('SystemUserData.pdf');
  };

  const handleSearchChange = (e) => {
    //console.log("Search Query:", e.target.value); // Debugging statement
    setSearchQuery(e.target.value); // Step 2: Update search query state
  };




  //	------For daterange---------//

  const dateRangePickerRef = useRef(null);

  useEffect(() => {
    $(document).ready(function () {
      // Log a message to check if this block is reached
      console.log('Document is ready');

      // Customize the Date Range Picker options
      const dateRangePickerOptions = {
    opens: 'left',
    autoUpdateInput: false,
    locale: {
      format: 'MM/DD/YYYY',
      applyLabel: t('Apply'),
      cancelLabel: t('Cancel'),
      customRangeLabel: t('Custom Range'),
    },
    ranges: {
      [t('Today')]: [moment(), moment()],
      [t('Yesterday')]: [moment().subtract(1, 'days'), moment().subtract(1, 'days')],
      [t('Last 7 Days')]: [moment().subtract(6, 'days'), moment()],
      [t('Last 30 Days')]: [moment().subtract(29, 'days'), moment()],
      [t('This Month')]: [moment().startOf('month'), moment().endOf('month')],
      [t('Last Month')]: [
        moment().subtract(1, 'month').startOf('month'),
        moment().subtract(1, 'month').endOf('month'),
      ],
        },
      };

      // Initialize the date range picker with options
      $(dateRangePickerRef.current).daterangepicker(dateRangePickerOptions);

      // Log a message to check if the Date Range Picker is initialized
      console.log('Date Range Picker is initialized');

      // Listen for date range changes
      $(dateRangePickerRef.current).on('apply.daterangepicker', (event, picker) => {
        const startDate = picker.startDate.format('MM/DD/YYYY');
        const endDate = picker.endDate.format('MM/DD/YYYY');

        // Log the selected date range
        console.log(`Selected date range: ${startDate} - ${endDate}`);

        // Set the selected date range to the state
        setCustomDateRange([startDate, endDate]);

        // Update the input field with the selected date range
        setCustomDateRangeText(`${startDate} - ${endDate}`);
      });
    });
  }, []);



  const [customDateRange, setCustomDateRange] = useState([
    lastMonthStart.format('MM/DD/YYYY'),
    today.format('MM/DD/YYYY'),
  ]);

  const linkStyle = {
    textDecoration: "none",
    color: '#808080'
  };
  const searchUserBreadcrumbs = [
    <Breadcrumbs aria-label="breadcrumb">
      <Link style={linkStyle} to="/Dashboard">
        <HomeIcon sx={{ mr: 0.5 }} fontSize="inherit" />
        {t('Home')}
      </Link>
      <Link  style={linkStyle} to="#">
      {t('System User Management')}
      </Link>
      <Typography color="black">{t('searchUser')}</Typography>
    </Breadcrumbs>
  ];


  //  --------For status dropdown------

  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(SEARCH_USER_URL);
        if (response.data && response.data.CustomerData) {
          setTableData(response.data.CustomerData);
          setRows(response.data.CustomerData);
        } else {
          setRows([]); // Set rows to an empty array if no data is available
        }
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    };

    fetchData();
  }, []);


  const [formData, setFormData] = useState({
    fullName: '',
    email: '',
    input_status: '',
  });


  const handleFormSubmit = async () => {
    try {
      const response = await axios.post(SEARCH_USER_URL, {
        ...formData,
        customStartDate: customDateRange[0],
        customEndDate: customDateRange[1],
      });

      if (response.data && response.data.CustomerDataPostMethod) {
        if (response.data.CustomerDataPostMethod.length === 0) {
          setRows([]); // Set rows to an empty array to trigger "No data available" message
          setTableData([]); // Update the tableData state with the new data
        } else {
          setRows(response.data.CustomerDataPostMethod); // Set rows with the response data
          setTableData(response.data.CustomerDataPostMethod); // Update the tableData state with the new data
        }
      } else {
        setRows([]); // Set rows to an empty array if the data is undefined or missing
        setTableData([]); // Update the tableData state with the new data
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error searching users:', error);
      setRows([]); // Set rows to an empty array on error to show "No data available"
      setTableData([]); // Update the tableData state with the new data
    }
  };

  const handleChange = (name, value) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };

  return (
    <>
      <MainCard
        title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('searchUser')}</span>}
        breadcrumbs={searchUserBreadcrumbs}
      >
        <Card sx={{ overflow: 'hidden' }}>
          <Grid container spacing={2}>
            <Grid item xs={3}>
              <TextField
                label={t('fullNameLabel')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.fullName}
                onChange={(e) => setFormData({ ...formData, fullName: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <TextField
                label={t('emailLabel')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <FormControl fullWidth variant="outlined" margin="normal">
                <InputLabel id="status-label" className={errors.language ? 'input-label-error' : ''}>
                  {t('statusLabel')}
                </InputLabel>
                <Select
                  labelId="status-label"
                  id="status"
                  value={formData.input_status}
                  onChange={(e) => handleChange('input_status', e.target.value)}
                  label="Status"
                  error={!!errors.status}
                >
                  <MenuItem value="">{t('selectStatus')}</MenuItem>
                  {hardcodedStatuses.map((status) => (
                    <MenuItem key={status.id} value={status.id}>
                      {status.name}
                    </MenuItem>
                  ))}
                </Select>
                {errors.status && <FormHelperText error>{errors.status}</FormHelperText>}
              </FormControl>
            </Grid>
            <Grid item xs={3}>
              <div className="form-group" style={{ position: 'relative' }}>
                <TextField
                  label={t('customDateRange')}
                  variant="outlined"
                  margin="normal"
                  value={customDateRangeText}
                  style={{width:'80%'}}
                  InputProps={{
                    readOnly: true,
                    onClick: () => {
                      // Log a message to check if the click event is triggered
                      console.log('Text field clicked'); 

                      // Show the date range picker when the text field is clicked
                      $(dateRangePickerRef.current).click();
                    },
                  }}
                />
                {/* DateRangePicker will be initialized here */}
                <div ref={dateRangePickerRef} style={{ position: 'absolute', top: '100%', left: 0 }} />
              </div>
            </Grid>
          </Grid>
          <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', p: 2 }}>
            <Button variant="contained" color="primary" style={{ backgroundColor: '#1478c8' }} onClick={handleFormSubmit}>
              {t('searchButton')}
            </Button>
            <Box sx={{ mx: 2 }} />
            <Button variant="contained" color="secondary" style={{ backgroundColor: 'rgb(13 32 97)' }} onClick={() => navigate('/SysUser/CreatesysUser')}>
              {t('createButton')}
            </Button>
          </Box>
        </Card>
      </MainCard>
      <br />
      <MainCard
        title=""
        secondary={<SecondaryAction icon={<LinkIcon fontSize="small" />} link="https://tablericons.com/" />}
      >
    <div>
        {/* copy button */}
        <Grid container spacing={2}>
          <Grid item xs={10}>
            <Button variant="contained" startIcon={<ContentCopyIcon />} className="exportbutton" onClick={() => setCopied(false)}>
              <CopyToClipboard text={JSON.stringify(tableData)} onCopy={() => setCopied(true)}>
                <span>{t('Copy')}</span>
              </CopyToClipboard>
            </Button>
            {copied ? <span style={{ color: 'green' }}>Copied!</span> : null}
            {/* copy button */}

            {/* excel button */}
            <Button variant="contained" startIcon={<InsertDriveFileIcon />} className="exportbutton" onClick={exportToExcel}>
              {t('Excel')}
            </Button>

            {/* pdf button */}
            <Button variant="contained" startIcon={<PictureAsPdfIcon />} className="exportbutton" onClick={exportToPDF}>
             {t('Pdf')}
            </Button>
          </Grid>



          <Grid item xs={2}>
            <TextField
              placeholder={t('searchButton')}
              variant="outlined"
              fullWidth
              margin="normal"
              style={{marginTop:'0px'}}
              value={searchQuery} // Use the search query to filter the table
              onChange={handleSearchChange}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <SearchIcon />
                  </InputAdornment>
                ),
                classes: {
                  notchedOutline: 'no-border', // Apply custom class to remove border
                },
              }}
            />
          </Grid>
        </Grid>

      </div>
        <Card sx={{ overflow: 'hidden' }}>{/* Your content here */}</Card>
        <TableContainer component={Paper}>

          <Table id="tableId">
            <TableHead style={{ backgroundColor: '#125e1f' }}>
              <TableRow>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('idLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('timestampLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('fullNameLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('mobileNumberLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('emailLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('addressLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('roleNameLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('language')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('statusLabel')}</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {(() => {

                // Filter rows based on searchQuery
                const filteredRows = rows.filter((row) => {
                  for (const key in row) {
                    if (typeof row[key] === 'string' && row[key].toLowerCase().includes(searchQuery.toLowerCase())) {
                      return true; // If a match is found in any property, include the row
                    }
                  }
                  return false; // If no match is found in any property, exclude the row
                });


                const startIndex = page * rowsPerPage;
                const endIndex = startIndex + rowsPerPage;

                if (rows.length === 0) {
                  return (
                    <TableRow>
                      <TableCell colSpan={8} align="center">
                        {t('noDataAvailable')}
                      </TableCell>
                    </TableRow>
                  );
                } else {
                  return filteredRows.slice(startIndex, endIndex).map((row) => (
                    <TableRow key={row.Id}>
                      <TableCell>
                        <Link to={`/SysUser/ViewsysUser/${row.Id}`} style={{ textDecoration: 'none' }}>
                          <Typography variant="h5" color="error">
                            {row.PartialId}
                          </Typography>
                        </Link>
                      </TableCell>
                      <TableCell>{row.Timestamp}</TableCell>
                      <TableCell>{row.FullName}</TableCell>
                      <TableCell>{row.Mobile}</TableCell>
                      <TableCell>{row.Email}</TableCell>
                      <TableCell>{row.Address}</TableCell>
                      <TableCell>{row.Rolename}</TableCell>
                      <TableCell>{row.Language}</TableCell>
                      <TableCell>{row.Status}</TableCell>
                    </TableRow>
                  ));
                }
              })()}
            </TableBody>

          </Table>

          <TablePagination
            rowsPerPageOptions={[5, 10, 25, 50]}
            component="div"
            count={rows.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={(event, newPage) => setPage(newPage)}
            onRowsPerPageChange={(event) => {
              setRowsPerPage(parseInt(event.target.value, 10));
              setPage(0);
            }}
            labelRowsPerPage={t('rowsPerPageLabel')}
            labelDisplayedRows={({ from, to, count }) => `${from}-${to} of ${count}`}
            // Custom style for the pagination controls
            className="pagination-container"
          />
        </TableContainer>
      </MainCard>
    </>
  );
};

export default TablerIcons;