import React, { useState, useEffect, useRef } from 'react';
import { styled } from '@mui/material/styles';
import { useNavigate } from 'react-router-dom';
import { TablePagination } from '@mui/material';
import { useTranslation } from 'react-i18next';
import { Table, TableBody, TableCell, TableContainer, TableHead, TableRow, Paper, Typography } from '@mui/material';
import { Card, Grid, Box, Button, TextField, FormControl, FormHelperText, InputLabel, Select, MenuItem, IconButton, InputAdornment } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import InsertDriveFileIcon from '@mui/icons-material/InsertDriveFile';
import PictureAsPdfIcon from '@mui/icons-material/PictureAsPdf';
import SearchIcon from '@mui/icons-material/Search';


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


//for making HTTP requests
import axios from 'axios';
import { Link } from 'react-router-dom';

// project imports
import MainCard from '../../../ui-component/cards/MainCard';
import SecondaryAction from '../../../ui-component/cards/CardSecondaryAction';
import { CREATE_USER_URL, SEARCH_USER_URL, ESBLOGS_REPORT_URL } from '../../../UrlPath.js';
import Breadcrumbs from '@mui/material/Breadcrumbs';
import HomeIcon from '@mui/icons-material/Home';

// assets
import LinkIcon from '@mui/icons-material/Link';
import '../../../assets/scss/style.css'
import Modal from '@mui/material/Modal';

// styles
const IFrameWrapper = styled('iframe')(({ theme }) => ({
  height: 'calc(100vh - 210px)',
  border: '1px solid',
  borderColor: theme.palette.primary.light
}));

// =============================|| TABLER ICONS ||============================= //

const ESBLogsReport = () => {
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
  const [rowsPerPage, setRowsPerPage] = useState(5);

  // Data to copy or export to xlsx or pdf in table
  const [copied, setCopied] = useState(false);
  const [tableData, setTableData] = useState([]);
  const [searchQuery, setSearchQuery] = useState('');
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [modalData, setModalData] = useState('');

  // Add the exportToExcel function
  const exportToExcel = () => {
    // Customize the data before exporting
    const customizedTableData = tableData.map((row) => ({
      'Id': row.Id,
      'Timestamp': row.Timestamp,
      'RequestId': row.RequestId,
      'Url': row.Url,
      'Service': row.Service,
      'In_Request': row.In_Request,
      'Out_Response': row.Out_Response,
      'ProducerAccessCode': row.ProducerAccessCode,

    }));

    const ws = XLSX.utils.json_to_sheet(customizedTableData);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, 'Sheet1');
    XLSX.writeFile(wb, 'ESBLogsReportData.xlsx');
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
      <Link style={linkStyle} to="#">
        {t('Reports')}
      </Link>
      <Typography color="black">{t('ESBLogsReport')}</Typography>
    </Breadcrumbs>
  ];



  const [customDateRange, setCustomDateRange] = useState([
    lastMonthStart.format('MM/DD/YYYY'),
    today.format('MM/DD/YYYY'),
  ]);
  //  --------For status dropdown------

  const hardcodedStatuses = [
    { id: 'ACTIVE', name: 'Active' },
    { id: 'INACTIVE', name: 'Inactive' },
  ];

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response = await axios.get(ESBLOGS_REPORT_URL);

        console.log(response.data.esb_logs_data);

        //to solve null data length error

        if (response.data && response.data.esb_logs_data) {
          setTableData(response.data.esb_logs_data);
          setRows(response.data.esb_logs_data);
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
    roleName: '',
    input_status: '',
  });


  const handleFormSubmit = async () => {
    try {
      const response = await axios.post(ESBLOGS_REPORT_URL, {
        ...formData,
        customStartDate: customDateRange[0],
        customEndDate: customDateRange[1],
      });

      if (response.data && response.data.esblogdata) {
        if (response.data.esblogdata.length === 0) {
          setRows([]);
          setTableData([]); // Set rows to an empty array to trigger "No data available" message
        } else {
          setRows(response.data.esblogdata);
          setTableData(response.data.esblogdata);  // Set rows with the response data
        }
      } else {
        setRows([]);
        setTableData([]); // Set rows to an empty array if the data is undefined or missing
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error searching users:', error);
      setRows([]);
      setTableData([]); // Set rows to an empty array on error to show "No data available"
    }
  };

  const handleChange = (name, value) => {
    setFormData((prevFormData) => ({
      ...prevFormData,
      [name]: value,
    }));
  };
  const handleShowPopup = (jsonData) => {
    setModalData(JSON.parse(jsonData)); // Parse the JSON string to an object
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };


  return (
    <>
      <MainCard
        title={<span style={{ color: 'black ', fontWeight: 'bold', fontSize: '24px' }}>{t('ESBLogsReport')}</span>}
        // style={{ backgroundColor: '#FF7F50' }}
        breadcrumbs={searchUserBreadcrumbs}
      >
        <Card sx={{ overflow: 'hidden' }}>
          <Grid container spacing={2} style={{ marginLeft: '10px' }}>
            <Grid item xs={3}>
              <TextField
                label={t('Service')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.service}
                onChange={(e) => setFormData({ ...formData, service: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <TextField
                label={t('RequestID')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.requestid}
                onChange={(e) => setFormData({ ...formData, requestid: e.target.value })}
              />
            </Grid>
            <Grid item xs={3}>
              <TextField
                label={t('Url')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.Url}
                onChange={(e) => setFormData({ ...formData, Url: e.target.value })}
              />
            </Grid>
            {/* <Grid item xs={3}>
              <TextField
                label={t('ProducerAccesscode')}
                variant="outlined"
                fullWidth
                margin="normal"
                value={formData.produceraccesscode}
                onChange={(e) => setFormData({ ...formData, produceraccesscode: e.target.value })}
              />
            </Grid> */}
            {/* <Grid item xs={3}>
              <FormControl fullWidth variant="outlined" margin="normal">
                <InputLabel id="status-label" className={errors.language ? 'input-label-error' : ''}>
                  Status
                </InputLabel>
                <Select
                  labelId="status-label"
                  id="status"
                  value={formData.input_status}
                  onChange={(e) => handleChange('input_status', e.target.value)}
                  label="Status"
                  error={!!errors.status}
                >
                  <MenuItem value="">Select Status</MenuItem>
                  {hardcodedStatuses.map((status) => (
                    <MenuItem key={status.id} value={status.id}>
                      {status.name}
                    </MenuItem>
                  ))}
                </Select>
                {errors.status && <FormHelperText error>{errors.status}</FormHelperText>}
              </FormControl>
            </Grid> */}
            <Grid item xs={3}>
              <div className="form-group" style={{ position: 'relative' }}>
                <TextField
                  label={t('customDateRange')}
                  variant="outlined"
                  margin="normal"
                  value={customDateRangeText}
                  style={{ width:'80%' }}
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
            <Button variant="contained" color="secondary" style={{ backgroundColor: '#1478c8' }} onClick={handleFormSubmit}>
              {t('searchButton')}
            </Button>
            <Box sx={{ mx: 2 }} />

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
                {/*<TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('idLabel')}</TableCell>*/}
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('timestampLabel')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('RequestID')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('Url')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('Service')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('InRequest')}</TableCell>
                <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('OutResponse')}</TableCell>
                {/* <TableCell style={{ color: 'white', fontWeight: 'bold' }}>{t('ProducerAccesscode')}</TableCell> */}
              </TableRow>
            </TableHead>
            <TableBody>
              {(() => {

                // Filter rows based on searchQuery
                const filteredRows = (rows || []).filter((row) => {
                  for (const key in row) {
                    if (row[key] && typeof row[key] === 'string' && row[key].toLowerCase().includes(searchQuery.toLowerCase())) {
                      return true; // If a match is found in any property, include the row
                    }
                  }
                  return false; // If no match is found in any property, exclude the row
                });

                const startIndex = page * rowsPerPage;
                const endIndex = startIndex + rowsPerPage;

                if (!filteredRows.length) {
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
                      {/*<TableCell>{row.Id}</TableCell>*/}
                      <TableCell>{row.Timestamp}</TableCell>
                      <TableCell>{row.RequestId}</TableCell>
                      <TableCell>{row.Url}</TableCell>
                      <TableCell>{row.Service}</TableCell>
                      <TableCell>
                        <Button variant="outlined" onClick={() => handleShowPopup(row.In_Request)}>
                        {t('showData')}
                        </Button>
                      </TableCell>
                      <TableCell>
                        <Button variant="outlined" onClick={() => handleShowPopup(row.Out_Response)}>
                        {t('showData')}
                        </Button>
                      </TableCell>

                      {/* <TableCell>{row.ProducerAccessCode}</TableCell> */}
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
         {/* Modal to display JSON data */}
      {modalData && (
        <Modal open={isModalOpen} onClose={handleCloseModal}>
          <Box
            sx={{
              position: 'absolute',
              top: '50%',
              left: '50%',
              transform: 'translate(-50%, -50%)',
              width: '60%',
              bgcolor: 'background.paper',
              border: '2px solid #000',
              boxShadow: 24,
              p: 2,
            }}
          >
            <Box
              sx={{
                maxHeight: '400px',
                overflow: 'auto',
                whiteSpace: 'pre-wrap',
                fontFamily: 'monospace',
              }}
            >
              {JSON.stringify(modalData, null, 2)}
            </Box>
            <Button onClick={handleCloseModal} variant="outlined">
            {t('close')}
            </Button>
          </Box>
        </Modal>
      )}
      </MainCard>
    </>
  );
};

export default ESBLogsReport;