UPDATE csv_imports
SET status=$3, reason=$4, data=$5
WHERE id=$1 AND user_id=$2