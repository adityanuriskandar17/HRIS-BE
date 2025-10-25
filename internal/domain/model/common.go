package model

type UserRole string

const (
	RoleEmployee UserRole = "EMPLOYEE"
	RoleManager  UserRole = "MANAGER"
	RoleHR       UserRole = "HR"
	RoleFinance  UserRole = "FINANCE"
	RoleAdmin    UserRole = "ADMIN"
)

type AttendanceMethod string

const (
	MethodWeb    AttendanceMethod = "WEB"
	MethodMobile AttendanceMethod = "MOBILE"
	MethodFace   AttendanceMethod = "FACE"
	MethodGeo    AttendanceMethod = "GEO"
)

type LeaveType string

const (
	LeaveAnnual LeaveType = "ANNUAL"
	LeaveSick   LeaveType = "SICK"
	LeaveUnpaid LeaveType = "UNPAID"
	LeaveOther  LeaveType = "OTHER"
)

type LeaveStatus string

const (
	LeavePending   LeaveStatus = "PENDING"
	LeaveApproved  LeaveStatus = "APPROVED"
	LeaveRejected  LeaveStatus = "REJECTED"
	LeaveCancelled LeaveStatus = "CANCELLED"
)

type EmploymentStatus string

const (
	EmploymentFullTime EmploymentStatus = "FULLTIME"
	EmploymentContract EmploymentStatus = "CONTRACT"
	EmploymentIntern   EmploymentStatus = "INTERN"
	EmploymentPartTime EmploymentStatus = "PARTTIME"
)
