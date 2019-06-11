import {h, app} from '/static/3rd/ha.js'

const Unlock = (state, action) => [{...state, lock: state.lock - 1}, action]
const Lock = (state, action) => [{...state, lock: state.lock + 1}, action]
const UpdateUsersUnlock = (state, users) => Unlock({...state, users})
const UpdateEmailsUnlock = (state, emails) =>
    Unlock({...state, emails: [...state.emails, ...emails]})
const UpdateNewUser = (state, e) => ({...state, newUser: e.target.value})
const AddNewEmail = (state, user) => ({
    ...state,
    newEmail: {user, address: ''},
})
const UpdateNewEmail = (state, e) => ({
    ...state,
    newEmail: {...state.newEmail, address: e.target.value},
})

const _fetchFx = ({path, action, options}, dispatch) =>
    fetch(path, options)
        .then(data => {
            return data.json()
        })
        .then(data => {
            dispatch(action, data)
        })

const FetchEmails = (state, id) => [
    Lock,
    [_fetchFx, {action: UpdateEmailsUnlock, path: '/rest/emails/' + id}],
]

const FetchUsers = state => [
    Lock,
    [_fetchFx, {action: UpdateUsersUnlock, path: '/rest/users'}],
]

const AddUser = state => [
    Lock,
    [
        _fetchFx,
        {
            action: [Unlock],
            path: '/rest/users/add',
            options: {
                method: 'POST',
                body: JSON.stringify({name: state.newUser}),
            },
        },
    ],
]

const AddEmail = state => [
    Lock,
    [
        _fetchFx,
        {
            action: [Unlock],
            path: '/rest/emails/add',
            options: {
                method: 'POST',
                body: JSON.stringify(state.newEmail),
            },
        },
    ],
]

app({
    init: () => ({
        lock: 0,
        newUser: '',
        newEmail: {},
        emails: [],
        users: [],
    }),
    view: state =>
        h('div', {}, [
            h('h1', {}, 'Rest'),
            h('a', {href: '/graphql/play'}, 'graphql'),
            h(
                'div',
                null,
                h(
                    'button',
                    {onclick: FetchUsers, disabled: state.lock},
                    'fetch users'
                )
            ),
            h(
                'div',
                null,
                h('input', {
                    type: 'text',
                    disabled: state.lock,
                    oninput: UpdateNewUser,
                    value: state.newUser,
                }),
                h('button', {disabled: state.lock, onclick: AddUser}, 'add user')
            ),
            h(
                'ul',
                null,
                state.users.map(u =>
                    h(
                        'li',
                        null,
                        h(
                            'strong',
                            {
                                'data-id': u.id,
                                'onclick': [FetchEmails, u.id],
                            },
                            u.name
                        ),
                        h(
                            'button',
                            {disabled: state.lock, onclick: [AddNewEmail, u.id]},
                            'add email'
                        ),
                        state.newEmail &&
              state.newEmail.user == u.id &&
              h(
                  'div',
                  null,
                  h('input', {
                      type: 'text',
                      disabled: state.lock,
                      value: state.newEmail.address,
                      oninput: UpdateNewEmail,
                  }),
                  h('button', {disabled: state.lock, onclick: AddEmail}, 'save')
              ),
                        state.emails
                            .filter(e => e.user_id == u.id)
                            .map(e => h('ol', {}, e.address))
                    )
                )
            ),
        ]),
    node: document.getElementById('app'),
})
