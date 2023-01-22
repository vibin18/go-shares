package dbrepo

import (
	"context"
	"github.com/vibin18/go-shares/internal/models"
	"time"
)

func (m *postgresDBRepo) InsertNewShare(res models.Share) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into share_names (name, id, created_at, updated_at) values ($1, $2, $3, $4)`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetShareByID(id int) (models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var share models.Share
	query := `select Name, id, created_at, updated_at from share_names where id = $2`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&share.Name,
		&share.Id,
		&share.CreatedAt,
		&share.UpdatedAt,
	)
	if err != nil {
		return share, err
	}
	return share, nil
}

func (m *postgresDBRepo) GetAllShares() ([]models.Share, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.Share
	query := `select name, id from share_names`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.Share
		err := rows.Scan(
			&r.Name,
			&r.Id,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}
func (m *postgresDBRepo) BuyShare(res models.SellBuyShare) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into purchased_shares (name, id, count, price, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		res.Count,
		res.Price,
		res.CreatedAt,
		res.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) SellShare(res models.SellBuyShare) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into sold_shares (name, id, count, price, created_at, updated_at) values ($1, $2, $3, $4, $5, $6)`
	_, err := m.DB.ExecContext(ctx, stmt,
		res.Name,
		res.Id,
		res.Count,
		res.Price,
		res.CreatedAt,
		res.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetAllSharesWithData() ([]models.TotalShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.TotalShare
	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.id, purchased_shares.name, purchased_shares.count as pc, coalesce(sold_shares.count, 0) as sc from purchased_shares left outer join sold_shares on purchased_shares.id = sold_shares.id`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.TotalShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.PCount,
			&r.SCount,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (m *postgresDBRepo) GetAllPurchases() ([]models.SellBuyShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.SellBuyShare

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.id, purchased_shares.name, purchased_shares.count, purchased_shares.price, purchased_shares.updated_at, purchased_shares.created_at  from purchased_shares order by created_at`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.SellBuyShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.Count,
			&r.Price,
			&r.UpdatedAt,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (m *postgresDBRepo) GetAllSales() ([]models.SellBuyShare, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.SellBuyShare

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select sold_shares.id, sold_shares.name, sold_shares.count, sold_shares.price, sold_shares.updated_at, sold_shares.created_at  from sold_shares order by created_at`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.SellBuyShare
		err := rows.Scan(
			&r.Id,
			&r.Name,
			&r.Count,
			&r.Price,
			&r.UpdatedAt,
			&r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (m *postgresDBRepo) GetAllSalesReport() ([]models.ShareReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.ShareReport

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select sold_shares.name,sum(sold_shares.count) as count_sum, sum(sold_shares.count * sold_shares.price) as total from sold_shares group by sold_shares.name order by name`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.ShareReport
		err := rows.Scan(
			&r.Name,
			&r.Count,
			&r.Total,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}

func (m *postgresDBRepo) GetAllPurchaseReport() ([]models.ShareReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var shares []models.ShareReport

	//	only works if both tables has values
	//	query := `select purchased_shares.id, purchased_shares.count as pc, sold_shares.count as sc from purchased_shares inner join sold_shares on purchased_shares.id = sold_shares.id`

	query := `select purchased_shares.name, sum(purchased_shares.count) as pc, sum(purchased_shares.count * purchased_shares.price) as pprice from purchased_shares group by purchased_shares.name order by name `

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var r models.ShareReport
		err := rows.Scan(
			&r.Name,
			&r.Count,
			&r.Total,
		)
		if err != nil {
			return nil, err
		}
		shares = append(shares, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return shares, nil

}
