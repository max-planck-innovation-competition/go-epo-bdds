package epo_bbds

// AuthHeader is the HTTP header for the Authorization Token
const AuthHeader = "Authorization"

// EpoLoginEndpoint is the endpoint for the EPO login
const EpoLoginEndpoint = "https://login.epo.org/oauth2/aus3up3nz0N133c0V417/v1/token"

// EpoBddsFileEndpoint is the endpoint for the docdb frontfiles bucket
// GET https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/%s/delivery/%s/file/%s/download
var EpoBddsFileEndpoint = "https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/%s/delivery/%d/file/%d/download"

// EpoBddsProductEndpoint is the endpoint for the doc db product
var EpoBddsProductEndpoint = "https://publication-bdds.apps.epo.org/bdds/bdds-bff-service/prod/api/products/%s"

// EpoBddsBProductID is the product id for epo bulk datasets
type EpoBddsBProductID string

// EpoFullTextFrontFilesProductID is the EP full-text data - front file
const EpoFullTextFrontFilesProductID EpoBddsBProductID = "4"

// EpoDocDBFrontFilesProductID is the product id for the doc db
const EpoDocDBFrontFilesProductID EpoBddsBProductID = "3"

// EpoDocDBBackFilesProductID is the product id for the doc db back files
const EpoDocDBBackFilesProductID EpoBddsBProductID = "14"

// EpoPatstatGlobalProductID is the product id for the PATSTAT global
const EpoPatstatGlobalProductID EpoBddsBProductID = "17"

// EpoPatstatEpRegisterProductID is the product id for the PATSTAT ep register
const EpoPatstatEpRegisterProductID EpoBddsBProductID = "18"
