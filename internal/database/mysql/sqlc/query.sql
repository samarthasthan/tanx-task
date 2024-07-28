-- name: CreateAccount :exec
INSERT INTO Users (UserID, Name, Email, Password)
VALUES (?,?,?,?);


-- name: CreateVerification :exec
INSERT INTO Verifications (VerificationId, UserID, OTP, ExpiresAt)
VALUES (?,?,?,?);

-- name: GetUserIDByEmail :one
SELECT UserID FROM Users WHERE Email = ?;

-- name: GetOTP :one
SELECT OTP, ExpiresAt FROM Verifications WHERE UserID = ?;

-- name: VerifyAccount :exec
UPDATE Users SET IsVerified = 1 WHERE UserID = ?;

-- name: DeleteVerification :exec
DELETE FROM Verifications WHERE UserID = ?;

-- name: GetPasswordByEmail :one
SELECT Password FROM Users WHERE Email = ?;

-- name: GetUserByEmail :one
SELECT UserID, Name, Email, Password FROM Users WHERE Email = ?;

-- name: CreateAlert :exec
INSERT INTO Alerts (AlertID, UserID, Curreny, Price) VALUES (?,?,?,?);

-- name: DeleteAlert :exec
UPDATE Alerts SET Status = 'deleted' WHERE AlertID = ?;

-- name: GetAlertsByUser :many
SELECT AlertID, UserID, Curreny, Price, Status, CreatedAt, UpdatedAt FROM Alerts WHERE UserID = ?;

-- name: GetAlerts :many
SELECT a.AlertID, a.UserID, a.Curreny, a.Price, a.Status, a.CreatedAt, a.UpdatedAt, u.Email
FROM Alerts a
JOIN Users u ON a.UserID = u.UserID
WHERE a.Status = 'created';

-- name: UpdateAlertStatus :exec
UPDATE Alerts SET Status = ? WHERE AlertID = ?;