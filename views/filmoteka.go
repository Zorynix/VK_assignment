package views

import (
	"net/http"
)

func (view *View) ActorAddView() error {
	data, err := view.PG.ActorAdd(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) ActorEditView() error {
	data, err := view.PG.ActorEdit(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) ActorListView() error {
	data, err := view.PG.ActorList(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) ActorDeleteView() error {
	data, err := view.PG.ActorDelete(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieAddView() error {
	data, err := view.PG.MovieAdd(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieEditView() error {
	data, err := view.PG.MovieEdit(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieFindView() error {
	data, err := view.PG.MovieFind(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieListView() error {
	data, err := view.PG.MovieList(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieDeleteView() error {
	data, err := view.PG.MovieDelete(view.W, view.R)
	if err != nil {
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}
