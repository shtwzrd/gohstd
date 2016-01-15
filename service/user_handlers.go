package service

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	gohst "github.com/warreq/gohstd/common"
	bcrypt "golang.org/x/crypto/bcrypt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	var u gohst.User
	err := ParseJsonEntity(r, &u)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	err2 := u.Validate()
	if err2 != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	s, err3 := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err3 != nil {
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	err4 := userRepo.InsertUser(u, s)
	if err4 != nil {
		if strings.Contains(err.Error(), gohst.UserExistsError) {
			HttpError(w, http.StatusConflict,
				errors.New(fmt.Sprintf("%s: %s", err, u.Username)))
			return
		}
		HttpError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// UserLogin is a do-nothing endpoint that merely confirms that authentication
// was successful. To that end, it is necessary that the endpoint is wrapped in
// actual authentication middleware.
func UserLogin(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
}

func UserUploadProfilePicture(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	user := vars["username"]

	src, hdr, err := r.FormFile("photo")
	if err != nil {
		if err == http.ErrMissingFile || err == http.ErrNotMultipart {
			HttpError(w, http.StatusBadRequest, err)
			return
		}
	}
	defer src.Close()

	contentType := hdr.Header["Content-Type"][0]
	if contentType != "image/png" && contentType != "image/jpg" {
		HttpError(w, http.StatusBadRequest, errors.New("File must be png or jpg"))
		return
	}

	filename := fmt.Sprintf("%s-%d", user, time.Now().Unix())
	err = os.MkdirAll("./images", 0666)
	if err != nil {
		HttpError(w, http.StatusInternalServerError, errors.New(""))
		return
	}

	if contentType == "image/png" {
		filename += ".png"
		_, err = png.Decode(src)
	} else {
		filename += ".jpg"
		_, err = jpeg.Decode(src)
	}

	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	src.Seek(0, 0)
	dst, err := os.Create(path.Join("./images", filename))

	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}

	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		HttpError(w, http.StatusBadRequest, err)
		return
	}
	userRepo.UpdateUserPicture(user, path.Join("./images", filename))

}

func UserShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagShow(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagRename(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}

func UserTagDelete(w http.ResponseWriter, r *http.Request) {
	panic("Not yet implemented")
}
